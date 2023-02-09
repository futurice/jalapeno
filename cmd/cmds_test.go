package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/ory/dockertest"
	"github.com/spf13/cobra"
)

type projectDirectoryPathCtxKey struct{}
type recipesDirectoryPathCtxKey struct{}
type certDirectoryPathCtxKey struct{}
type ociRegistryCtxKey struct{}
type cmdStdOutCtxKey struct{}
type cmdStdErrCtxKey struct{}
type dockerResourcesCtxKey struct{}

type OCIRegistry struct {
	TLS      bool
	Auth     bool
	Resource *dockertest.Resource
}

/*
 * UTILITIES
 */

func WrapCmdOutputs(cmdFactory func() *cobra.Command) (*cobra.Command, *bytes.Buffer, *bytes.Buffer) {
	cmd := cmdFactory()
	cmdStdOut, cmdStdErr := new(bytes.Buffer), new(bytes.Buffer)
	cmd.SetOut(cmdStdOut)
	cmd.SetErr(cmdStdErr)

	return cmd, cmdStdOut, cmdStdErr
}

/*
 * STEP DEFINITIONS
 */

func aProjectDirectory(ctx context.Context) (context.Context, error) {
	dir, err := os.MkdirTemp("", "jalapeno-test-project")
	if err != nil {
		return ctx, err
	}

	return context.WithValue(ctx, projectDirectoryPathCtxKey{}, dir), nil
}

func aRecipesDirectory(ctx context.Context) (context.Context, error) {
	dir, err := os.MkdirTemp("", "jalapeno-test-recipes")
	if err != nil {
		return ctx, err
	}

	return context.WithValue(ctx, recipesDirectoryPathCtxKey{}, dir), nil
}

func aRecipeThatGeneratesFile(ctx context.Context, recipe, filename string) (context.Context, error) {
	dir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	if err := os.MkdirAll(filepath.Join(dir, recipe, "templates"), 0755); err != nil {
		return ctx, err
	}
	template := "apiVersion: v1\nname: %[1]s\nversion: v0.0.1\ndescription: %[1]s"
	if err := os.WriteFile(filepath.Join(dir, recipe, "recipe.yml"), []byte(fmt.Sprintf(template, recipe)), 0644); err != nil {
		return ctx, err
	}
	if err := os.WriteFile(filepath.Join(dir, recipe, "templates", filename), []byte(recipe), 0644); err != nil {
		return ctx, err
	}
	return context.WithValue(ctx, recipesDirectoryPathCtxKey{}, dir), nil
}

func aLocalOCIRegistry(ctx context.Context) (context.Context, error) {
	resource, err := createLocalRegistry(&dockertest.RunOptions{Repository: "registry", Tag: "2"})
	if err != nil {
		return ctx, err
	}

	ctx = context.WithValue(ctx, ociRegistryCtxKey{}, OCIRegistry{Resource: resource})
	ctx = addDockerResourceToContext(ctx, resource)

	return ctx, nil
}

func aLocalOCIRegistryWithAuth(ctx context.Context) (context.Context, error) {
	ctx, err := generateTLSCertificate(ctx)
	if err != nil {
		return ctx, err
	}

	resource, err := createLocalRegistry(&dockertest.RunOptions{
		Repository: "registry",
		Tag:        "2",
		Env: []string{
			// "REGISTRY_AUTH_SILLY_REALM=test-realm",
			// "REGISTRY_AUTH_SILLY_SERVICE=test-service",
			"REGISTRY_HTTP_TLS_CERTIFICATE=/etc/ssl/private/cert.pem",
			"REGISTRY_HTTP_TLS_KEY=/etc/ssl/private/key.pem",
		},
		Mounts: []string{fmt.Sprintf("%s:/etc/ssl/private", ctx.Value(certDirectoryPathCtxKey{}).(string))},
	})
	if err != nil {
		return ctx, err
	}

	ctx = context.WithValue(ctx, ociRegistryCtxKey{}, OCIRegistry{TLS: true, Auth: true, Resource: resource})
	ctx = addDockerResourceToContext(ctx, resource)

	return ctx, nil
}

func iExecuteRecipe(ctx context.Context, recipe string) (context.Context, error) {
	projectDir := ctx.Value(projectDirectoryPathCtxKey{}).(string)
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)

	cmd, cmdStdOut, cmdStdErr := WrapCmdOutputs(newExecuteCmd)

	cmd.SetArgs([]string{filepath.Join(recipesDir, recipe)})
	if err := cmd.Flags().Set("output", projectDir); err != nil {
		return ctx, err
	}

	if err := cmd.Execute(); err != nil {
		return ctx, err
	}

	ctx = context.WithValue(ctx, cmdStdOutCtxKey{}, cmdStdOut.String())
	ctx = context.WithValue(ctx, cmdStdErrCtxKey{}, cmdStdErr.String())

	return ctx, nil
}

func pushRecipe(ctx context.Context, recipeName, repoName string) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	registry := ctx.Value(ociRegistryCtxKey{}).(OCIRegistry)

	cmd, cmdStdOut, cmdStdErr := WrapCmdOutputs(newPushCmd)

	cmd.SetArgs([]string{filepath.Join(recipesDir, recipeName), filepath.Join(registry.Resource.GetHostPort("5000/tcp"), repoName)})

	flags := cmd.Flags()
	if registry.TLS {
		if err := flags.Set("insecure", "true"); err != nil {
			return ctx, err
		}
	} else {
		if err := flags.Set("plain-http", "true"); err != nil {
			return ctx, err
		}
	}

	if err := cmd.Execute(); err != nil {
		return ctx, err
	}

	ctx = context.WithValue(ctx, cmdStdOutCtxKey{}, cmdStdOut.String())
	ctx = context.WithValue(ctx, cmdStdErrCtxKey{}, cmdStdErr.String())

	return ctx, nil
}

func iPullRecipe(ctx context.Context, recipeName, repoName string) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	registry := ctx.Value(ociRegistryCtxKey{}).(OCIRegistry)

	cmd, cmdStdOut, cmdStdErr := WrapCmdOutputs(newPullCmd)

	cmd.SetArgs([]string{filepath.Join(registry.Resource.GetHostPort("5000/tcp"), repoName)})
	flags := cmd.Flags()
	if err := flags.Set("output", recipesDir); err != nil {
		return ctx, err
	}
	if registry.TLS {
		// Allow self-signed certificates
		if err := flags.Set("insecure", "true"); err != nil {
			return ctx, err
		}
	} else {
		if err := flags.Set("plain-http", "true"); err != nil {
			return ctx, err
		}
	}

	if err := cmd.Execute(); err != nil {
		return ctx, err
	}

	ctx = context.WithValue(ctx, cmdStdOutCtxKey{}, cmdStdOut.String())
	ctx = context.WithValue(ctx, cmdStdErrCtxKey{}, cmdStdErr.String())

	return ctx, nil
}

func pushOfTheRecipeWasSuccessful(ctx context.Context) (context.Context, error) {
	cmdStdOut := ctx.Value(cmdStdOutCtxKey{}).(string)
	cmdStdErr := ctx.Value(cmdStdErrCtxKey{}).(string)

	if cmdStdErr != "" {
		return ctx, fmt.Errorf("stderr was not empty: %s", cmdStdErr)
	}

	if cmdStdOut == "" { // TODO: Check stdout when we have proper message from CMD
		return ctx, errors.New("stdout was empty")
	}

	return ctx, nil
}

func pullOfTheRecipeWasSuccessful(ctx context.Context) (context.Context, error) {
	cmdStdOut := ctx.Value(cmdStdOutCtxKey{}).(string)
	cmdStdErr := ctx.Value(cmdStdErrCtxKey{}).(string)

	if cmdStdErr != "" {
		return ctx, fmt.Errorf("stderr was not empty: %s", cmdStdErr)
	}

	if cmdStdOut == "" { // TODO: Check stdout when we have proper message from CMD
		return ctx, errors.New("stdout was empty")
	}

	return ctx, nil
}

func theRecipesDirectoryShouldContainRecipe(ctx context.Context, recipeName string) error {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	re, err := recipe.Load(filepath.Join(recipesDir, recipeName))
	if err != nil {
		return err
	}

	if re.Name != recipeName {
		return fmt.Errorf("recipe name was \"%s\", expected \"%s\"", re.Name, recipeName)
	}

	return nil
}

func theProjectDirectoryShouldContainFile(ctx context.Context, filename string) error {
	dir := ctx.Value(projectDirectoryPathCtxKey{}).(string)
	info, err := os.Stat(filepath.Join(dir, filename))
	if err == nil && !info.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", filename)
	}
	return err
}

func cleanTempDirs(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
	if dir := ctx.Value(projectDirectoryPathCtxKey{}); dir != nil {
		os.RemoveAll(dir.(string))
	}
	if dir := ctx.Value(recipesDirectoryPathCtxKey{}); dir != nil {
		os.RemoveAll(dir.(string))
	}
	if dir := ctx.Value(certDirectoryPathCtxKey{}); dir != nil {
		os.RemoveAll(dir.(string))
	}
	return ctx, err
}

func cleanDockerResources(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
	resources, ok := ctx.Value(dockerResourcesCtxKey{}).([]*dockertest.Resource)

	// Resource list was probably empty, skip
	if !ok {
		return ctx, err
	}

	for _, resource := range resources {
		err := resource.Close()
		if err != nil {
			return ctx, err
		}
	}
	return ctx, err
}

func executionOfTheRecipeHasSucceeded(ctx context.Context) (context.Context, error) {
	cmdStdOut := ctx.Value(cmdStdOutCtxKey{}).(string)
	cmdStdErr := ctx.Value(cmdStdErrCtxKey{}).(string)
	if matched, _ := regexp.MatchString("Recipe executed successfully", cmdStdOut); !matched {
		return ctx, fmt.Errorf("Recipe failed to execute!\nstdout:\n%s\n\nstderr:\n%s\n", cmdStdOut, cmdStdErr)
	}
	return ctx, nil
}

func executionOfTheRecipeHasFailedWithError(ctx context.Context, errorMessage string) (context.Context, error) {
	cmdStdOut := ctx.Value(cmdStdOutCtxKey{}).(string)
	cmdStdErr := ctx.Value(cmdStdErrCtxKey{}).(string)
	if matched, _ := regexp.MatchString(errorMessage, cmdStdErr); !matched {
		return ctx, fmt.Errorf("'%s' not found in stderr.\nstdout:\n%s\n\nstderr:\n%s\n", errorMessage, cmdStdOut, cmdStdErr)
	}
	return ctx, nil
}

func iChangeRecipeToVersion(ctx context.Context, recipeName, version string) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	recipeFile := filepath.Join(recipesDir, recipeName, "recipe.yml")
	recipeData, err := os.ReadFile(recipeFile)
	if err != nil {
		return ctx, err
	}

	newData := strings.Replace(string(recipeData), "v0.0.1", version, 1)

	if err := os.WriteFile(filepath.Join(recipesDir, recipeName, "recipe.yml"), []byte(newData), 0644); err != nil {
		return ctx, err
	}

	return ctx, nil
}

func iUpgradeRecipe(ctx context.Context, recipe string) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	projectDir := ctx.Value(projectDirectoryPathCtxKey{}).(string)

	cmd, cmdStdOut, cmdStdErr := WrapCmdOutputs(newUpgradeCmd)

	cmd.SetArgs([]string{projectDir, filepath.Join(recipesDir, recipe)})

	if err := cmd.Execute(); err != nil {
		return ctx, err
	}

	ctx = context.WithValue(ctx, cmdStdOutCtxKey{}, cmdStdOut.String())
	ctx = context.WithValue(ctx, cmdStdErrCtxKey{}, cmdStdErr.String())

	return ctx, nil
}

func theProjectDirectoryShouldContainFileWith(ctx context.Context, filename, searchTerm string) error {
	cmdStdOut := ctx.Value(cmdStdOutCtxKey{}).(string)
	cmdStdErr := ctx.Value(cmdStdErrCtxKey{}).(string)
	dir := ctx.Value(projectDirectoryPathCtxKey{}).(string)

	path := filepath.Join(dir, filename)
	info, err := os.Stat(path)
	if err != nil {
		return err
	} else if !info.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file.\nstdout:\n%s\n\nstderr:\n%s\n", filename, cmdStdOut, cmdStdErr)
	}
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if !strings.Contains(string(bytes), searchTerm) {
		return fmt.Errorf("substring %s not found in %s.\nstdout:\n%s\n\nstderr:\n%s\n", searchTerm, filename, cmdStdOut, cmdStdErr)
	}
	return nil
}

func recipeIgnoresPattern(ctx context.Context, recipeName, pattern string) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	recipeFile := filepath.Join(recipesDir, recipeName, "recipe.yml")
	recipeData, err := os.ReadFile(recipeFile)
	if err != nil {
		return ctx, err
	}
	recipe := fmt.Sprintf("%s\nignorePatterns:\n  - %s\n", string(recipeData), pattern)
	if err := os.WriteFile(recipeFile, []byte(recipe), 0644); err != nil {
		return ctx, err
	}
	return ctx, nil
}

func iChangeProjectFileToContain(ctx context.Context, filename, content string) (context.Context, error) {
	projectDir := ctx.Value(projectDirectoryPathCtxKey{}).(string)
	if err := os.WriteFile(filepath.Join(projectDir, filename), []byte(content), 0644); err != nil {
		return ctx, err
	}
	return ctx, nil
}

func noConflictsWereReported(ctx context.Context) (context.Context, error) {
	cmdStdOut := ctx.Value(cmdStdOutCtxKey{}).(string)
	cmdStdErr := ctx.Value(cmdStdErrCtxKey{}).(string)
	if matched, _ := regexp.MatchString("modified", cmdStdOut); matched {
		return ctx, fmt.Errorf("Conflict in recipe\nstdout:\n%s\n\nstderr:\n%s\n", cmdStdOut, cmdStdErr)
	}
	return ctx, nil
}

func conflictsAreReported(ctx context.Context) (context.Context, error) {
	cmdStdOut := ctx.Value(cmdStdOutCtxKey{}).(string)
	cmdStdErr := ctx.Value(cmdStdErrCtxKey{}).(string)
	if matched, _ := regexp.MatchString("modified", cmdStdOut); matched {
		return ctx, nil
	}
	return ctx, fmt.Errorf("Expecting conflicts in recipe but none reported\nstdout:\n%s\n\nstderr:\n%s\n", cmdStdOut, cmdStdErr)
}

func iChangeRecipeTemplateToRender(ctx context.Context, recipeName, filename, content string) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	templateFilePath := filepath.Join(recipesDir, recipeName, "templates", filename)
	if err := os.WriteFile(templateFilePath, []byte(content), 0644); err != nil {
		return ctx, err
	}
	return ctx, nil
}

func noErrorsWerePrinted(ctx context.Context) (context.Context, error) {
	cmdStdErr := ctx.Value(cmdStdErrCtxKey{}).(string)
	if len(cmdStdErr) != 0 {
		return ctx, fmt.Errorf("Expected stderr to be empty but was %s", cmdStdErr)
	}
	return ctx, nil
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			s.Step(`^a project directory$`, aProjectDirectory)
			s.Step(`^a recipes directory$`, aRecipesDirectory)
			s.Step(`^a recipe "([^"]*)" that generates file "([^"]*)"$`, aRecipeThatGeneratesFile)
			s.Step(`^I execute recipe "([^"]*)"$`, iExecuteRecipe)
			s.Step(`^the project directory should contain file "([^"]*)"$`, theProjectDirectoryShouldContainFile)
			s.Step(`^the project directory should contain file "([^"]*)" with "([^"]*)"$`, theProjectDirectoryShouldContainFileWith)
			s.Step(`^execution of the recipe has succeeded$`, executionOfTheRecipeHasSucceeded)
			s.Step(`^execution of the recipe has failed with error "([^"]*)"$`, executionOfTheRecipeHasFailedWithError)
			s.Step(`^I change recipe "([^"]*)" to version "([^"]*)"$`, iChangeRecipeToVersion)
			s.Step(`^I upgrade recipe "([^"]*)"$`, iUpgradeRecipe)
			s.Step(`^recipe "([^"]*)" ignores pattern "([^"]*)"$`, recipeIgnoresPattern)
			s.Step(`^I change project file "([^"]*)" to contain "([^"]*)"$`, iChangeProjectFileToContain)
			s.Step(`^no conflicts were reported$`, noConflictsWereReported)
			s.Step(`^conflicts are reported$`, conflictsAreReported)
			s.Step(`^I change recipe "([^"]*)" template "([^"]*)" to render "([^"]*)"$`, iChangeRecipeTemplateToRender)
			s.Step(`^no errors were printed$`, noErrorsWerePrinted)
			s.Step(`^a local OCI registry$`, aLocalOCIRegistry)
			s.Step(`^a local OCI registry with authentication$`, aLocalOCIRegistryWithAuth)
			s.Step(`^I push the recipe "([^"]*)" to the local OCI repository "([^"]*)"$`, pushRecipe)
			s.Step(`^I pull the recipe "([^"]*)" to the local OCI repository "([^"]*)"$`, iPullRecipe)
			s.Step(`^the recipe "([^"]*)" is pushed to the local OCI repository "([^"]*)"$`, pushRecipe)
			s.Step(`^push of the recipe was successful$`, pushOfTheRecipeWasSuccessful)
			s.Step(`^pull of the recipe was successful$`, pullOfTheRecipeWasSuccessful)
			s.Step(`^the recipes directory should contain recipe "([^"]*)"$`, theRecipesDirectoryShouldContainRecipe)
			s.After(cleanDockerResources)
			s.After(cleanTempDirs)
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func TestExampleRecipe(t *testing.T) {
	recipe := createExampleRecipe("foo")
	if err := recipe.Validate(); err != nil {
		t.Errorf("failed to validate the example recipe: %s", err)
	}
}

// UTILS

func createLocalRegistry(opts *dockertest.RunOptions) (*dockertest.Resource, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, fmt.Errorf("could not construct pool: %w", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		return nil, fmt.Errorf("could not connect to Docker: %s", err)
	}

	resource, err := pool.RunWithOptions(opts)
	if err != nil {
		return nil, fmt.Errorf("could not start resource: %s", err)
	}

	host := resource.GetHostPort("5000/tcp")

	pool.MaxWait = 30 * time.Second
	if err = pool.Retry(func() error {
		_, err := pool.Client.HTTPClient.Get(fmt.Sprintf("http://%s/v2/", host))
		return err
	}); err != nil {
		return nil, fmt.Errorf("could not connect to docker: %s", err)
	}

	// Even though we check if the registry is ready, running tests immediately causes EOF errors to happen.
	// So we need to wait a bit more to registry to be ready.
	time.Sleep(100 * time.Millisecond)

	resource.Expire(60) // If the cleanup fails, this will stop the container eventually

	return resource, nil
}

func addDockerResourceToContext(ctx context.Context, resource *dockertest.Resource) context.Context {
	resources, ok := ctx.Value(dockerResourcesCtxKey{}).([]*dockertest.Resource)
	if !ok {
		resources = make([]*dockertest.Resource, 0)
	}

	return context.WithValue(ctx, dockerResourcesCtxKey{}, append(resources, resource))
}

func generateTLSCertificate(ctx context.Context) (context.Context, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return ctx, err
	}
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 180),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return ctx, err
	}

	dir, err := os.MkdirTemp("", "jalapeno-test-certs")
	if err != nil {
		return ctx, err
	}

	cert := &bytes.Buffer{}
	pem.Encode(cert, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	err = os.WriteFile(filepath.Join(dir, "cert.pem"), cert.Bytes(), 0666)
	if err != nil {
		return ctx, err
	}

	key := &bytes.Buffer{}
	b, _ := x509.MarshalECPrivateKey(priv)
	pem.Encode(key, &pem.Block{Type: "EC PRIVATE KEY", Bytes: b})
	err = os.WriteFile(filepath.Join(dir, "key.pem"), key.Bytes(), 0666)
	if err != nil {
		return ctx, err
	}

	return context.WithValue(ctx, certDirectoryPathCtxKey{}, dir), nil
}
