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
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/cucumber/godog"
	"github.com/gofrs/uuid"
	"github.com/muesli/termenv"
	"github.com/ory/dockertest/v3"

	"github.com/futurice/jalapeno/internal/cli"
	"github.com/futurice/jalapeno/pkg/recipe"
)

type (
	projectDirectoryPathCtxKey      struct{}
	recipesDirectoryPathCtxKey      struct{}
	certDirectoryPathCtxKey         struct{}
	htpasswdDirectoryPathCtxKey     struct{}
	dockerConfigDirectoryPathCtxKey struct{}
	manifestDirectoryPathCtxKey     struct{}
	ociRegistryCtxKey               struct{}
	cmdStdOutCtxKey                 struct{}
	cmdStdInCtxKey                  struct{}
	cmdStdErrCtxKey                 struct{}
	cmdAdditionalFlagsCtxKey        struct{}
	dockerResourcesCtxKey           struct{}
	scenarioNameCtxKey              struct{}
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
		Options: &godog.Options{
			Format:      "pretty",
			Strict:      true,
			Concurrency: runtime.NumCPU(),
			Paths:       []string{"features"},
			TestingT:    t,
		},
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			AddCommonSteps(s)

			// Command specific steps
			AddBumpverSteps(s)
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

			// Disable colors when testing
			s.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
				lipgloss.SetColorProfile(termenv.Ascii)
				return ctx, nil
			})

			// Initialize context values
			s.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
				ctx = context.WithValue(ctx, cmdAdditionalFlagsCtxKey{}, make(map[string]string))
				ctx = context.WithValue(ctx, dockerResourcesCtxKey{}, []*dockertest.Resource{})
				ctx = context.WithValue(ctx, cmdStdInCtxKey{}, NewBlockBuffer())
				ctx = context.WithValue(ctx, scenarioNameCtxKey{}, sc.Name)

				dirs := map[interface{}]string{
					projectDirectoryPathCtxKey{}:  "jalapeno-test-project",
					recipesDirectoryPathCtxKey{}:  "jalapeno-test-recipes",
					manifestDirectoryPathCtxKey{}: "jalapeno-test-manifest",
				}

				for key, dirPrefix := range dirs {
					dir, err := os.MkdirTemp("", dirPrefix)
					if err != nil {
						return ctx, err
					}
					ctx = context.WithValue(ctx, key, dir)
				}

				return ctx, nil
			})

			s.After(cleanDockerResources)
			s.After(cleanTempDirs)
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func AddCommonSteps(s *godog.ScenarioContext) {
	// Setup steps
	s.Step(`^a recipe "([^"]*)"$`, aRecipe)
	s.Step(`^recipe "([^"]*)" generates file "([^"]*)" with content "([^"]*)"$`, recipeGeneratesFileWithContent)
	s.Step(`^recipe "([^"]*)" ignores pattern "([^"]*)"$`, recipeIgnoresPattern)
	s.Step(`^I remove file "([^"]*)" from the recipe "([^"]*)"$`, iRemoveFileFromTheRecipe)
	s.Step(`^I create a file "([^"]*)" with contents "([^"]*)" to the project directory$`, iCreateAFileWithContentsToTheProjectDir)
	s.Step(`^a local OCI registry$`, aLocalOCIRegistry)
	s.Step(`^a local OCI registry with authentication$`, aLocalOCIRegistryWithAuth)
	s.Step(`^registry credentials are not provided by the command$`, credentialsAreNotProvidedByTheCommand)
	s.Step(`^registry credentials are provided by config file$`, generateDockerConfigFile)
	s.Step(`^I buffer key presses "([^"]*)"$`, bufferKeysToInput)
	s.Step(`^I clear the output$`, iClearTheOutput)

	// Assert steps
	s.Step(`^the recipes directory should contain recipe "([^"]*)"$`, theRecipesDirectoryShouldContainRecipe)
	s.Step(`^the manifest directory should contain manifest named "([^"]*)"$`, theManifestDirectoryShouldContainManifest)
	s.Step(`^no errors were printed$`, noErrorsWerePrinted)
	s.Step(`^CLI produced an output "([^"]*)"$`, expectGivenOutput)
	s.Step(`^CLI produced an error "(.*)"$`, expectGivenError)
	s.Step(`^the sauce in index (\d) which should have property "([^"]*)"$`, theSauceShouldHaveProperty)
	s.Step(`^the sauce in index (\d) which should not have property "([^"]*)"$`, theSauceFileShouldNotHaveProperty)
	s.Step(`^the sauce in index (\d) which should have property "([^"]*)" with value "([^"]*)"$`, theSauceFileShouldHavePropertyWithValue)
	s.Step(`^the sauce in index (\d) which has a valid ID$`, theSauceFileShouldHasAValidID)
	s.Step(`^the project directory should contain file "([^"]*)"$`, theProjectDirectoryShouldContainFile)
	s.Step(`^the project directory should contain file "([^"]*)" with "([^"]*)"$`, theProjectDirectoryShouldContainFileWith)
	s.Step(`^the project directory should not contain file "([^"]*)"$`, theProjectDirectoryShouldNotContainFile)
	s.Step(`^the file "([^"]*)" exist in the recipe "([^"]*)"$`, theFileExistInTheRecipe)
}

/*
 * UTILITIES
 */

func executeCLI(ctx context.Context, args ...string) (context.Context, error) {
	rootCmd := cli.NewRootCmd()

	cmdStdOut, cmdStdErr := new(bytes.Buffer), new(bytes.Buffer)
	cmdStdIn := ctx.Value(cmdStdInCtxKey{}).(*BlockBuffer)
	ctx = context.WithValue(ctx, cmdStdOutCtxKey{}, cmdStdOut)
	ctx = context.WithValue(ctx, cmdStdErrCtxKey{}, cmdStdErr)

	rootCmd.SetOut(cmdStdOut)
	rootCmd.SetIn(cmdStdIn)
	rootCmd.SetErr(cmdStdErr)
	rootCmd.SetContext(context.Background())

	additionalFlags := ctx.Value(cmdAdditionalFlagsCtxKey{}).(map[string]string)
	for name, value := range additionalFlags {
		args = append(args, fmt.Sprintf("--%s=%s", name, value))
	}

	rootCmd.SetArgs(args)
	_ = rootCmd.Execute()

	ctx = clearAdditionalFlags(ctx)
	return ctx, nil
}

func cleanTempDirs(ctx context.Context, sc *godog.Scenario, lastStepErr error) (context.Context, error) {
	directoryCtxKeys := []interface{}{
		projectDirectoryPathCtxKey{},
		recipesDirectoryPathCtxKey{},
		certDirectoryPathCtxKey{},
		htpasswdDirectoryPathCtxKey{},
		dockerConfigDirectoryPathCtxKey{},
		manifestDirectoryPathCtxKey{},
	}

	for _, key := range directoryCtxKeys {
		if dir := ctx.Value(key); dir != nil {
			os.RemoveAll(dir.(string))
		}
	}

	return ctx, nil
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

func readSauces(ctx context.Context) ([]*recipe.Sauce, error) {
	dir := ctx.Value(projectDirectoryPathCtxKey{}).(string)

	sauces, err := recipe.LoadSauces(dir)
	if err != nil {
		return nil, err
	}

	return sauces, nil
}

/*
 * STEP DEFINITIONS
 */

func cleanDockerResources(ctx context.Context, sc *godog.Scenario, lastStepErr error) (context.Context, error) {
	resources := ctx.Value(dockerResourcesCtxKey{}).([]*dockertest.Resource)

	for _, resource := range resources {
		err := resource.Close()
		if err != nil {
			return ctx, err
		}
	}
	return ctx, nil
}

func bufferKeysToInput(ctx context.Context, keys string) (context.Context, error) {
	stdIn := ctx.Value(cmdStdInCtxKey{}).(*BlockBuffer)

	commandChars := []string{
		"\\r", "\r",
		"\\n", "\n",
		"\\x1b", "\x1b",
	}

	// From https://github.com/charmbracelet/bubbletea/blob/master/key.go#L354
	customMappings := []string{
		"↑", "\x1b[A",
		"↓", "\x1b[B",
		"→", "\x1b[C",
		"←", "\x1b[D",
	}

	customKeys := make([]string, 0, len(customMappings)/2)
	for i := 0; i < len(customMappings); i += 2 {
		customKeys = append(customKeys, customMappings[i])
	}

	r := regexp.MustCompile(fmt.Sprintf("(%s)", strings.Join(customKeys, "|")))
	splitters := r.FindAllString(keys, -1)
	other := r.Split(keys, -1)

	blocks := []string{other[0]}
	for i := range splitters {
		blocks = append(blocks, splitters[i], other[i+1])
	}

	replacer := strings.NewReplacer(append(commandChars, customMappings...)...)
	for _, block := range blocks {
		if block != "" {
			stdIn.AddBlock([]byte(replacer.Replace(block)))
		}
	}

	return context.WithValue(ctx, cmdStdInCtxKey{}, stdIn), nil
}

func aRecipe(ctx context.Context, recipeName string) (context.Context, error) {
	dir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	scenarioName := ctx.Value(scenarioNameCtxKey{}).(string)

	re := recipe.NewRecipe()
	re.Name = recipeName
	re.Version = "v0.0.1"
	re.Description = scenarioName

	if err := re.Save(filepath.Join(dir, recipeName)); err != nil {
		return ctx, err
	}

	return ctx, nil
}

func recipeGeneratesFileWithContent(ctx context.Context, recipeName, filename, content string) (context.Context, error) {
	dir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	re, err := recipe.LoadRecipe(filepath.Join(dir, recipeName))
	if err != nil {
		return ctx, err
	}

	re.Templates[filename] = recipe.NewFile([]byte(content))

	if err := re.Save(filepath.Join(dir, recipeName)); err != nil {
		return ctx, err
	}

	return ctx, nil
}

func iRemoveFileFromTheRecipe(ctx context.Context, filename, recipeName string) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)

	templateDir := filepath.Join(recipesDir, recipeName, recipe.RecipeTemplatesDirName)

	err := os.Remove(filepath.Join(templateDir, filename))
	return ctx, err
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
	re, err := recipe.LoadRecipe(filepath.Join(recipesDir, recipeName))
	if err != nil {
		return err
	}

	if re.Name != recipeName {
		return fmt.Errorf("recipe name was \"%s\", expected \"%s\"", re.Name, recipeName)
	}

	return nil
}

func theManifestDirectoryShouldContainManifest(ctx context.Context, manifestName string) error {
	manifestDir := ctx.Value(manifestDirectoryPathCtxKey{}).(string)
	info, err := os.Stat(filepath.Join(manifestDir, manifestName))
	if err == nil && !info.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", manifestName)
	}

	return err
}

func theProjectDirectoryShouldNotContainFile(ctx context.Context, filename string) error {
	err := theProjectDirectoryShouldContainFile(ctx, filename)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	return err
}

func iClearTheOutput(ctx context.Context) (context.Context, error) {
	ctx = context.WithValue(ctx, cmdStdOutCtxKey{}, new(bytes.Buffer))
	ctx = context.WithValue(ctx, cmdStdErrCtxKey{}, new(bytes.Buffer))
	return ctx, nil
}

func theProjectDirectoryShouldContainFile(ctx context.Context, filename string) error {
	dir := ctx.Value(projectDirectoryPathCtxKey{}).(string)
	info, err := os.Stat(filepath.Join(dir, filename))
	if err == nil && !info.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", filename)
	}
	return err
}

func iCreateAFileWithContentsToTheProjectDir(ctx context.Context, filename, contents string) error {
	dir := ctx.Value(projectDirectoryPathCtxKey{}).(string)
	return os.WriteFile(filepath.Join(dir, filename), []byte(contents), 0644)
}

func theProjectDirectoryShouldContainFileWith(ctx context.Context, filename, searchTerm string) error {
	content, err := readProjectDirectoryFile(ctx, filename)
	if err != nil {
		return err
	}

	if matched, err := regexp.MatchString(searchTerm, content); err != nil {
		return fmt.Errorf("regexp pattern matching caused an error: %w", err)
	} else if !matched {
		return fmt.Errorf("the file '%s' did not match the following pattern '%s'", filename, searchTerm)
	}

	return nil
}

func theSauceFileShouldHasAValidID(ctx context.Context, index int) error {
	sauces, err := readSauces(ctx)
	if err != nil {
		return err
	}

	if sauces[index].ID == uuid.Nil {
		return errors.New("recipe file does not have 'id' property")
	}

	return nil
}

func theSauceFileShouldHavePropertyWithValue(ctx context.Context, index int, propertyName, expectedValue string) error {
	sauces, err := readSauces(ctx)
	if err != nil {
		return err
	}

	r := getDeepPropertyFromStruct(sauces[index], propertyName)
	if !r.IsValid() {
		return fmt.Errorf("sauce file does not have property %s", propertyName)
	}

	value := r.String()

	if matched, err := regexp.MatchString(expectedValue, value); err != nil {
		return fmt.Errorf("regexp pattern matching caused an error: %w", err)
	} else if !matched {
		return fmt.Errorf("expected property %s to match regex '%s', got '%s'", propertyName, expectedValue, value)
	}

	return nil
}

func theSauceShouldHaveProperty(ctx context.Context, index int, propertyName string) error {
	sauces, err := readSauces(ctx)
	if err != nil {
		return err
	}

	r := getDeepPropertyFromStruct(sauces[index], propertyName)
	if !r.IsValid() {
		return fmt.Errorf("sauce file does not have property %s", propertyName)
	}

	return nil
}

func theSauceFileShouldNotHaveProperty(ctx context.Context, index int, propertyName string) error {
	err := theSauceShouldHaveProperty(ctx, index, propertyName)
	if err == nil {
		return fmt.Errorf("sauce file contains the property %s: %w", propertyName, err)
	}

	if !strings.Contains(err.Error(), "not have property") {
		return err
	}

	return nil
}

func recipeIgnoresPattern(ctx context.Context, recipeName, pattern string) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	re, err := recipe.LoadRecipe(filepath.Join(recipesDir, recipeName))
	if err != nil {
		return ctx, err
	}

	re.IgnorePatterns = append(re.IgnorePatterns, pattern)

	if err := re.Save(filepath.Join(recipesDir, recipeName)); err != nil {
		return ctx, err
	}

	return ctx, nil
}

func expectGivenOutput(ctx context.Context, expected string) error {
	cmdStdOut := ctx.Value(cmdStdOutCtxKey{}).(*bytes.Buffer)

	if matched, err := regexp.MatchString(expected, cmdStdOut.String()); err != nil {
		return fmt.Errorf("regexp pattern matching caused an error: %w", err)
	} else if !matched {
		return fmt.Errorf("command produced unexpected output: Expected: '%s', Actual: '%s'", expected, strings.TrimSpace(cmdStdOut.String()))
	}

	return nil
}

func expectGivenError(ctx context.Context, expectedError string) error {
	cmdStdErr := ctx.Value(cmdStdErrCtxKey{}).(*bytes.Buffer)

	if matched, err := regexp.MatchString(expectedError, cmdStdErr.String()); err != nil {
		return fmt.Errorf("regexp pattern matching caused an error: %w", err)
	} else if !matched {
		return fmt.Errorf("command produced unexpected error: Expected: '%s', Actual: '%s'", expectedError, strings.TrimSpace(cmdStdErr.String()))
	}

	return nil
}

func noErrorsWerePrinted(ctx context.Context) error {
	cmdStdErr := ctx.Value(cmdStdErrCtxKey{}).(*bytes.Buffer)
	if cmdStdErr.String() != "" {
		return fmt.Errorf("Expected stderr to be empty but was '%s'", strings.TrimSpace(cmdStdErr.String()))
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
	resources := ctx.Value(dockerResourcesCtxKey{}).([]*dockertest.Resource)
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

func addRegistryRelatedFlags(ctx context.Context) context.Context {
	ociRegistry := ctx.Value(ociRegistryCtxKey{}).(OCIRegistry)
	configDir, configFileExists := ctx.Value(dockerConfigDirectoryPathCtxKey{}).(string)
	additionalFlags := ctx.Value(cmdAdditionalFlagsCtxKey{}).(map[string]string)

	if ociRegistry.TLSEnabled {
		// Allow self-signed certificates
		additionalFlags["insecure"] = "true"
	} else {
		additionalFlags["plain-http"] = "true"
	}

	if ociRegistry.AuthEnabled {
		additionalFlags["username"] = "foo"
		additionalFlags["password"] = "bar"
	}

	if configFileExists && os.Getenv("DOCKER_CONFIG") == "" {
		additionalFlags["registry-config"] = filepath.Join(configDir, DOCKER_CONFIG_FILENAME)
	}

	return ctx
}

func getDeepPropertyFromStruct(v any, key string) reflect.Value {
	r := reflect.ValueOf(v)
	for _, k := range strings.Split(key, ".") {
		switch reflect.Indirect(r).Kind() {
		case reflect.Struct:
			r = reflect.Indirect(r).FieldByName(k)
		case reflect.Map:
			r = r.MapIndex(reflect.ValueOf(k))
		default:
			return r
		}
	}

	return r
}

func clearAdditionalFlags(ctx context.Context) context.Context {
	return context.WithValue(ctx, cmdAdditionalFlagsCtxKey{}, make(map[string]string))
}

// BlockBuffer represents a buffer that stores blocks of data.
// This is used to simulate user input for the CLI. Using the standard buffer
// would not allow the simulation of complex key presses due to how bubbletea is implemented.
type BlockBuffer struct {
	data       [][]byte // The blocks of data stored in the buffer.
	blockIndex int      // The index of the current block.
	readIndex  int      // The index of the next byte to be read.
}

var _ io.Reader = &BlockBuffer{}

func NewBlockBuffer() *BlockBuffer {
	return &BlockBuffer{
		data:       make([][]byte, 0),
		blockIndex: 0,
		readIndex:  0,
	}
}

func (r *BlockBuffer) AddBlock(p []byte) {
	r.data = append(r.data, p)
}

func (r *BlockBuffer) Len() int {
	s := 0
	for i := range r.data {
		s += len(r.data[i])
	}

	return s
}

func (r *BlockBuffer) Read(p []byte) (n int, err error) {
	if r.blockIndex >= len(r.data) || r.readIndex >= len(r.data[r.blockIndex]) {
		err = io.EOF
		return
	}

	n = copy(p, r.data[r.blockIndex][r.readIndex:])
	if n == len(r.data[r.blockIndex]) {
		r.blockIndex++
		r.readIndex = 0
	} else {
		r.readIndex += n
	}
	return
}
