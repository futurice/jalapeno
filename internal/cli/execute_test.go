package cli_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"
)

func AddExecuteSteps(s *godog.ScenarioContext) {
	s.Step(`^I execute recipe "([^"]*)"$`, iRunExecute)
	s.Step(`^I execute the recipe from the local OCI repository "([^"]*)"$`, iExecuteRemoteRecipe)
	s.Step(`^recipes will be executed to the subpath "([^"]*)"$`, recipesWillBeExecutedToTheSubPath)
	s.Step(`^execution of the recipe has succeeded$`, executionOfTheRecipeHasSucceeded)
}

func iRunExecute(ctx context.Context, recipe string) (context.Context, error) {
	projectDir := ctx.Value(projectDirectoryPathCtxKey{}).(string)
	stdIn := ctx.Value(cmdStdInCtxKey{}).(*BlockBuffer)

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

	if stdIn.Len() == 0 {
		args = append(args, "--no-input")
	}

	return executeCLI(ctx, args...)
}

func iExecuteRemoteRecipe(ctx context.Context, repository string) (context.Context, error) {
	ociRegistry := ctx.Value(ociRegistryCtxKey{}).(OCIRegistry)
	configDir, configFileExists := ctx.Value(dockerConfigDirectoryPathCtxKey{}).(string)
	additionalFlags := ctx.Value(cmdAdditionalFlagsCtxKey{}).(map[string]string)

	url := fmt.Sprintf("oci://%s/%s", ociRegistry.Resource.GetHostPort("5000/tcp"), repository)

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

	return iRunExecute(ctx, url)
}

func recipesWillBeExecutedToTheSubPath(ctx context.Context, path string) (context.Context, error) {
	additionalFlags := ctx.Value(cmdAdditionalFlagsCtxKey{}).(map[string]string)
	additionalFlags["subpath"] = path

	return ctx, nil
}

func executionOfTheRecipeHasSucceeded(ctx context.Context) (context.Context, error) {
	return ctx, expectGivenOutput(ctx, "Recipe executed successfully")
}
