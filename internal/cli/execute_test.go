package cli_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"
	re "github.com/futurice/jalapeno/pkg/recipe"
)

func AddExecuteSteps(s *godog.ScenarioContext) {
	s.Step(`^I execute recipe "([^"]*)"$`, iRunExecute)
	s.Step(`^I execute the recipe from the local OCI repository "([^"]*)"$`, iExecuteRemoteRecipe)
	s.Step(`^execution of the recipe has succeeded$`, executionOfTheRecipeHasSucceeded)
	s.Step(`^a manifest file that includes recipes "([^"]*)" and "([^"]*)"$`, aManifestFileThatIncludesRecipesAnd)
	s.Step(`^I execute the manifest file$`, iExecuteTheManifestFile)
}

func iRunExecute(ctx context.Context, recipe string) (context.Context, error) {
	projectDir := ctx.Value(projectDirectoryPathCtxKey{}).(string)
	additionalFlags := ctx.Value(cmdAdditionalFlagsCtxKey{}).(map[string]string)

	ctx, cmd := wrapCmdOutputs(ctx)

	var url string
	if strings.HasPrefix(recipe, "oci://") {
		url = recipe
	} else {
		recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
		url = filepath.Join(recipesDir, recipe)
	}

	args := []string{
		"execute",
		url,
		fmt.Sprintf("--dir=%s", projectDir),
	}

	for name, value := range additionalFlags {
		args = append(args, fmt.Sprintf("--%s=%s", name, value))
	}

	cmd.SetArgs(args)
	_ = cmd.Execute()
	return ctx, nil
}

func iExecuteRemoteRecipe(ctx context.Context, repository string) (context.Context, error) {
	ociRegistry := ctx.Value(ociRegistryCtxKey{}).(OCIRegistry)
	configDir, configFileExists := ctx.Value(dockerConfigDirectoryPathCtxKey{}).(string)
	flags := ctx.Value(cmdAdditionalFlagsCtxKey{}).(map[string]string)

	url := fmt.Sprintf("oci://%s/%s", ociRegistry.Resource.GetHostPort("5000/tcp"), repository)

	if ociRegistry.TLSEnabled {
		// Allow self-signed certificates
		flags["insecure"] = "true"
	} else {
		flags["plain-http"] = "true"
	}

	if ociRegistry.AuthEnabled {
		flags["username"] = "foo"
		flags["password"] = "bar"
	}

	if configFileExists && os.Getenv("DOCKER_CONFIG") == "" {
		flags["registry-config"] = filepath.Join(configDir, DOCKER_CONFIG_FILENAME)
	}

	ctx = context.WithValue(ctx, cmdAdditionalFlagsCtxKey{}, flags)

	return iRunExecute(ctx, url)
}

func executionOfTheRecipeHasSucceeded(ctx context.Context) (context.Context, error) {
	return ctx, expectGivenOutput(ctx, "Recipe executed successfully")
}

func aManifestFileThatIncludesRecipesAnd(ctx context.Context, recipe1, recipe2 string) (context.Context, error) {
	recipeDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	dir, err := os.MkdirTemp("", "jalapeno-test-manifest")
	if err != nil {
		return ctx, err
	}
	ctx = context.WithValue(ctx, manifestDirectoryPathCtxKey{}, dir)
	manifest := fmt.Sprintf(`apiVersion: v1
recipes:
  - name: %[2]s
    version: 0.0.0
    repository: file://%[1]s/%[2]s
  - name: %[3]s
    version: 0.0.0
    repository: file://%[1]s/%[3]s
`, recipeDir, recipe1, recipe2)
	if err := os.WriteFile(filepath.Join(dir, re.ManifestFileName+re.YAMLExtension), []byte(manifest), 0644); err != nil {
		return ctx, err
	}
	return ctx, godog.ErrPending
}

func iExecuteTheManifestFile() error {
	return godog.ErrPending
}
