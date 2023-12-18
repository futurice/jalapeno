package cli_test

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
	"io"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/gofrs/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/futurice/jalapeno/internal/cli"
	re "github.com/futurice/jalapeno/pkg/recipe"
)

type (
	projectDirectoryPathCtxKey      struct{}
	recipesDirectoryPathCtxKey      struct{}
	certDirectoryPathCtxKey         struct{}
	htpasswdDirectoryPathCtxKey     struct{}
	dockerConfigDirectoryPathCtxKey struct{}
	ociRegistryCtxKey               struct{}
	cmdStdOutCtxKey                 struct{}
	cmdStdErrCtxKey                 struct{}
	cmdAdditionalFlagsCtxKey        struct{}
	dockerResourcesCtxKey           struct{}
)

type OCIRegistry struct {
	TLSEnabled  bool
	AuthEnabled bool
	Resource    *dockertest.Resource
}

const (
	TLS_KEY_FILENAME         = "key.pem"
	TLS_CERTIFICATE_FILENAME = "cert.pem"
	HTPASSWD_FILENAME        = "htpasswd"
	DOCKER_CONFIG_FILENAME   = "config.json"
)

/*
 * STEP DEFINITIONS
 */

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			// Common steps
			s.Step(`^a project directory$`, aProjectDirectory)
			s.Step(`^a recipes directory$`, aRecipesDirectory)
			s.Step(`^a recipe "([^"]*)" that generates file "([^"]*)"$`, aRecipeThatGeneratesFile)
			s.Step(`^a failing recipe "([^"]*)" with variable "([^"]*)" that generates file "([^"]*)"$`, aFailingRecipeWithVariableThatGeneratesFile)
			s.Step(`^the file "([^"]*)" exist in the recipe "([^"]*)"$`, theFileExistInTheRecipe)
			s.Step(`^the project directory should contain file "([^"]*)"$`, theProjectDirectoryShouldContainFile)
			s.Step(`^the project directory should contain file "([^"]*)" with "([^"]*)"$`, theProjectDirectoryShouldContainFileWith)
			s.Step(`^the sauce file contains a sauce in index (\d) which should have property "([^"]*)" with value "([^"]*)"$`, theSauceFileShouldHavePropertyWithValue)
			s.Step(`^the sauce file contains a sauce in index (\d) which should have property "([^"]*)" that is a valid UUID$`, theSauceFileShouldHavePropertyThatIsAValidUUID)
			s.Step(`^CLI produced an output "([^"]*)"$`, expectGivenOutput)
			s.Step(`^CLI produced an error "([^"]*)"$`, expectGivenError)
			s.Step(`^recipe "([^"]*)" ignores pattern "([^"]*)"$`, recipeIgnoresPattern)
			s.Step(`^no errors were printed$`, noErrorsWerePrinted)
			s.Step(`^a local OCI registry$`, aLocalOCIRegistry)
			s.Step(`^a local OCI registry with authentication$`, aLocalOCIRegistryWithAuth)
			s.Step(`^registry credentials are not provided by the command$`, credentialsAreNotProvidedByTheCommand)
			s.Step(`^registry credentials are provided by config file$`, generateDockerConfigFile)
			s.Step(`^the recipes directory should contain recipe "([^"]*)"$`, theRecipesDirectoryShouldContainRecipe)

			// Command specific steps
			AddCheckSteps(s)
			AddCreateSteps(s)
			AddEjectSteps(s)
			AddExecuteSteps(s)
			AddPullSteps(s)
			AddPushSteps(s)
			AddTestSteps(s)
			AddUpgradeSteps(s)
			AddValidateSteps(s)
			AddWhySteps(s)

			s.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
				// Initialize additional flags to empty map before each step
				return context.WithValue(
					ctx,
					cmdAdditionalFlagsCtxKey{},
					make(map[string]string)), nil
			})
			s.After(cleanDockerResources)
			s.After(cleanTempDirs)
		},
		Options: &godog.Options{
			Strict:      true,
			Concurrency: 4,
			Format:      "pretty",
			Paths:       []string{"../../test"},
			TestingT:    t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

/*
 * UTILITIES
 */

func wrapCmdOutputs(ctx context.Context) (context.Context, *cobra.Command) {
	rootCmd := cli.NewRootCmd()
	cmdStdOut, cmdStdErr := new(bytes.Buffer), new(bytes.Buffer)

	rootCmd.SetOut(cmdStdOut)
	rootCmd.SetErr(cmdStdErr)
	rootCmd.SetContext(context.Background())

	ctx = context.WithValue(ctx, cmdStdOutCtxKey{}, cmdStdOut)
	ctx = context.WithValue(ctx, cmdStdErrCtxKey{}, cmdStdErr)

	return ctx, rootCmd
}

func cleanTempDirs(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
	directoryCtxKeys := []interface{}{
		projectDirectoryPathCtxKey{},
		recipesDirectoryPathCtxKey{},
		certDirectoryPathCtxKey{},
		htpasswdDirectoryPathCtxKey{},
		dockerConfigDirectoryPathCtxKey{},
	}

	for _, key := range directoryCtxKeys {
		if dir := ctx.Value(key); dir != nil {
			os.RemoveAll(dir.(string))
		}
	}

	return ctx, err
}

func readProjectDirectoryFile(ctx context.Context, filename string) (string, error) {
	dir := ctx.Value(projectDirectoryPathCtxKey{}).(string)
	path := filepath.Join(dir, filename)
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	} else if !info.Mode().IsRegular() {
		return "", fmt.Errorf("%s is not a regular file", filename)
	}
	bytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func readSauceFile(ctx context.Context) ([]map[interface{}]interface{}, error) {
	content, err := readProjectDirectoryFile(ctx, filepath.Join(re.SauceDirName, re.SaucesFileName+re.YAMLExtension))
	if err != nil {
		return nil, err
	}
	var recipes []map[interface{}]interface{}
	decoder := yaml.NewDecoder(bytes.NewReader([]byte(content)))
	for {
		recipe := make(map[interface{}]interface{})
		if err := decoder.Decode(&recipe); err != nil {
			if err != io.EOF {
				return nil, fmt.Errorf("failed to decode recipe file: %w", err)
			}
			break
		}
		recipes = append(recipes, recipe)
	}
	return recipes, nil
}

/*
 * STEP DEFINITIONS
 */

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
	template := `apiVersion: v1
name: %[1]s
version: v0.0.1
description: %[1]s
`
	if err := os.WriteFile(filepath.Join(dir, recipe, re.RecipeFileName+re.YAMLExtension), []byte(fmt.Sprintf(template, recipe)), 0644); err != nil {
		return ctx, err
	}

	templateDir := filepath.Join(dir, recipe, re.RecipeTemplatesDirName)
	err := os.MkdirAll(filepath.Join(templateDir, filepath.Dir(filename)), 0755)
	if err != nil {
		return ctx, err
	}

	if err := os.WriteFile(filepath.Join(templateDir, filename), []byte(recipe), 0644); err != nil {
		return ctx, err
	}
	return context.WithValue(ctx, recipesDirectoryPathCtxKey{}, dir), nil
}

func aFailingRecipeWithVariableThatGeneratesFile(ctx context.Context, recipe, variable, filename string) (context.Context, error) {
	dir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	if err := os.MkdirAll(filepath.Join(dir, recipe, "templates"), 0755); err != nil {
		return ctx, err
	}
	template := `apiVersion: v1
name: %[1]s
version: v0.0.1
description: %[1]s
variables:
  %[2]s:
    type: string
	description: %[2]s
`
	if err := os.WriteFile(filepath.Join(dir, recipe, re.RecipeFileName+re.YAMLExtension), []byte(fmt.Sprintf(template, recipe, variable)), 0644); err != nil {
		return ctx, err
	}

	templateDir := filepath.Join(dir, recipe, re.RecipeTemplatesDirName)
	err := os.MkdirAll(filepath.Join(templateDir, filepath.Dir(filename)), 0755)
	if err != nil {
		return ctx, err
	}

	if err := os.WriteFile(filepath.Join(templateDir, filename), []byte("{{ .failboat }}"), 0644); err != nil {
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

func credentialsAreNotProvidedByTheCommand(ctx context.Context) (context.Context, error) {
	registry := ctx.Value(ociRegistryCtxKey{}).(OCIRegistry)
	registry.AuthEnabled = false

	return context.WithValue(ctx, ociRegistryCtxKey{}, registry), nil
}

func theRecipesDirectoryShouldContainRecipe(ctx context.Context, recipeName string) error {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	re, err := re.LoadRecipe(filepath.Join(recipesDir, recipeName))
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
	content, err := readProjectDirectoryFile(ctx, filename)
	if err != nil {
		return err
	}
	if !strings.Contains(content, searchTerm) {
		return fmt.Errorf("substring '%s' not found in file %s", searchTerm, filename)
	}
	return nil
}

func theSauceFileShouldHavePropertyThatIsAValidUUID(ctx context.Context, index int, propertyName string) error {
	recipes, err := readSauceFile(ctx)
	if err != nil {
		return err
	}

	value, exists := (recipes[index])[propertyName].(string)
	if exists {
		if _, err := uuid.FromString(value); err != nil {
			return fmt.Errorf("found UUID but it does not parse: %w", err)
		}
	} else {
		return fmt.Errorf("recipe file does not have property %s", propertyName)
	}

	if err != nil {
		return err
	}

	return nil
}

func theSauceFileShouldHavePropertyWithValue(ctx context.Context, index int, propertyName, expectedValue string) error {
	recipes, err := readSauceFile(ctx)
	if err != nil {
		return err
	}
	value, exists := (recipes[index])[propertyName].(string)
	if !exists {
		return fmt.Errorf("recipe file does not have property %s", propertyName)
	}

	if !regexp.MustCompile(expectedValue).MatchString(value) {
		return fmt.Errorf("expected property %s to match regex '%s', got '%s'", propertyName, expectedValue, value)
	}
	return nil
}

func recipeIgnoresPattern(ctx context.Context, recipeName, pattern string) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	recipeFile := filepath.Join(recipesDir, recipeName, re.RecipeFileName+re.YAMLExtension)
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

func expectGivenOutput(ctx context.Context, expected string) error {
	cmdStdOut := ctx.Value(cmdStdOutCtxKey{}).(*bytes.Buffer)

	if matched, err := regexp.MatchString(expected, cmdStdOut.String()); !matched {
		return fmt.Errorf("command produced unexpected output: Expected: '%s', Actual: '%s'", expected, strings.TrimSpace(cmdStdOut.String()))
	} else if err != nil {
		return fmt.Errorf("regexp pattern matching caused an error: %w", err)
	}

	return nil
}

func expectGivenError(ctx context.Context, expectedError string) error {
	cmdStdErr := ctx.Value(cmdStdErrCtxKey{}).(*bytes.Buffer)

	if matched, err := regexp.MatchString(expectedError, cmdStdErr.String()); !matched {
		return fmt.Errorf("command produced unexpected error: Expected: '%s', Actual: '%s'", expectedError, strings.TrimSpace(cmdStdErr.String()))
	} else if err != nil {
		return fmt.Errorf("regexp pattern matching caused an error: %w", err)
	}

	return nil
}

func noErrorsWerePrinted(ctx context.Context) error {
	cmdStdErr := ctx.Value(cmdStdErrCtxKey{}).(*bytes.Buffer)
	if cmdStdErr.String() != "" {
		return fmt.Errorf("Expected stderr to be empty but was '%s'", cmdStdErr)
	}
	return nil
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

	pool.MaxWait = 10 * time.Second
	if err = pool.Retry(func() error {
		url := fmt.Sprintf("http://%s/v2/", host)
		resp, err := http.Get(url)
		if err != nil {
			return err
		}

		// Non-authenticated registry responds with status 200, authenticated with 400
		if resp.StatusCode == 200 || resp.StatusCode == 400 {
			return nil
		}

		return errors.New("endpoint not yet healthy")
	}); err != nil {
		return nil, fmt.Errorf("could not connect to docker: %s", err)
	}

	err = resource.Expire(60) // If the cleanup fails, this will stop the container eventually
	if err != nil {
		return nil, err
	}

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
	err = pem.Encode(cert, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		return ctx, err
	}

	err = os.WriteFile(filepath.Join(dir, TLS_CERTIFICATE_FILENAME), cert.Bytes(), 0666)
	if err != nil {
		return ctx, err
	}

	key := &bytes.Buffer{}
	b, _ := x509.MarshalECPrivateKey(priv)
	err = pem.Encode(key, &pem.Block{Type: "EC PRIVATE KEY", Bytes: b})
	if err != nil {
		return ctx, err
	}

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

	// Created with `docker run --entrypoint htpasswd httpd:2 -Bbn foo bar`
	contents := "foo:$2y$05$fHux.x9qjOuYmARV5AXPpuNnph95rssj5tsIeMynjL1O7jj43YMrW\n" // foo:bar
	err = os.WriteFile(filepath.Join(dir, HTPASSWD_FILENAME), []byte(contents), 0666)
	if err != nil {
		return ctx, err
	}

	return context.WithValue(ctx, htpasswdDirectoryPathCtxKey{}, dir), nil
}

func generateDockerConfigFile(ctx context.Context) (context.Context, error) {
	registry := ctx.Value(ociRegistryCtxKey{}).(OCIRegistry)
	dir, err := os.MkdirTemp("", "jalapeno-test-docker-config")
	if err != nil {
		return ctx, err
	}

	contents := fmt.Sprintf(`{
  "auths": {
    "https://%s/v2/": {
      "auth": "Zm9vOmJhcg=="
    }
  }
}`, registry.Resource.GetHostPort("5000/tcp"))
	err = os.WriteFile(filepath.Join(dir, DOCKER_CONFIG_FILENAME), []byte(contents), 0666)
	if err != nil {
		return ctx, err
	}

	return context.WithValue(ctx, dockerConfigDirectoryPathCtxKey{}, dir), nil
}

func theFileExistInTheRecipe(ctx context.Context, file, recipe string) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)

	path := filepath.Join(recipesDir, recipe, file)
	if info, err := os.Stat(path); os.IsNotExist(err) {
		return ctx, fmt.Errorf("the file %s does not exist", file)
	} else if info.IsDir() {
		return ctx, errors.New("the path contained a directory instead of a file")
	}

	return ctx, nil
}
