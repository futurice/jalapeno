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
	"fmt"
	"math/big"
	"os"
	"path/filepath"
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
type htpasswdDirectoryPathCtxKey struct{}
type ociRegistryCtxKey struct{}
type cmdStdOutCtxKey struct{}
type cmdStdErrCtxKey struct{}
type dockerResourcesCtxKey struct{}

type OCIRegistry struct {
	TLSEnabled  bool
	AuthEnabled bool
	Resource    *dockertest.Resource
}

const (
	TLS_KEY_FILENAME         = "key.pem"
	TLS_CERTIFICATE_FILENAME = "cert.pem"
	HTPASSWD_FILENAME        = "htpasswd"
)

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
			s.Step(`^I push the recipe "([^"]*)" to the local OCI repository "([^"]*)"$`, iPushRecipe)
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
	if dir := ctx.Value(htpasswdDirectoryPathCtxKey{}); dir != nil {
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

	ctx, err = generateHtpasswdFile(ctx)
	if err != nil {
		return ctx, err
	}

	resource, err := createLocalRegistry(&dockertest.RunOptions{
		Repository: "registry",
		Tag:        "2",
		Env: []string{
			"REGISTRY_AUTH_HTPASSWD_REALM=jalapeno-test-realm",
			fmt.Sprintf("REGISTRY_AUTH_HTPASSWD_PATH=/auth/%s", HTPASSWD_FILENAME),
			fmt.Sprintf("REGISTRY_HTTP_TLS_CERTIFICATE=/etc/ssl/private/%s", TLS_CERTIFICATE_FILENAME),
			fmt.Sprintf("REGISTRY_HTTP_TLS_KEY=/etc/ssl/private/%s", TLS_KEY_FILENAME),
		},
		Mounts: []string{
			fmt.Sprintf("%s:/etc/ssl/private", ctx.Value(certDirectoryPathCtxKey{}).(string)),
			fmt.Sprintf("%s:/auth", ctx.Value(htpasswdDirectoryPathCtxKey{}).(string)),
		},
	})
	if err != nil {
		return ctx, err
	}

	ctx = context.WithValue(ctx, ociRegistryCtxKey{}, OCIRegistry{TLSEnabled: true, AuthEnabled: true, Resource: resource})
	ctx = addDockerResourceToContext(ctx, resource)

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

func noErrorsWerePrinted(ctx context.Context) (context.Context, error) {
	cmdStdErr := ctx.Value(cmdStdErrCtxKey{}).(string)
	if len(cmdStdErr) != 0 {
		return ctx, fmt.Errorf("Expected stderr to be empty but was %s", cmdStdErr)
	}
	return ctx, nil
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
	err = os.WriteFile(filepath.Join(dir, TLS_CERTIFICATE_FILENAME), cert.Bytes(), 0666)
	if err != nil {
		return ctx, err
	}

	key := &bytes.Buffer{}
	b, _ := x509.MarshalECPrivateKey(priv)
	pem.Encode(key, &pem.Block{Type: "EC PRIVATE KEY", Bytes: b})
	err = os.WriteFile(filepath.Join(dir, TLS_KEY_FILENAME), key.Bytes(), 0666)
	if err != nil {
		return ctx, err
	}

	return context.WithValue(ctx, certDirectoryPathCtxKey{}, dir), nil
}

func generateHtpasswdFile(ctx context.Context) (context.Context, error) {
	dir, err := os.MkdirTemp("", "jalapeno-test-htpasswd")
	if err != nil {
		return ctx, err
	}

	// Username foo, password bar
	contents := "foo:$2y$05$fHux.x9qjOuYmARV5AXPpuNnph95rssj5tsIeMynjL1O7jj43YMrW\n"
	err = os.WriteFile(filepath.Join(dir, HTPASSWD_FILENAME), []byte(contents), 0666)
	if err != nil {
		return ctx, err
	}

	return context.WithValue(ctx, htpasswdDirectoryPathCtxKey{}, dir), nil
}
