package cli_test

import (
	"bytes"
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
	s.Step(`^execution of the recipe has succeeded$`, executionOfTheRecipeHasSucceeded)
}

func iRunExecute(ctx context.Context, recipe string) (context.Context, error) {
	projectDir := ctx.Value(projectDirectoryPathCtxKey{}).(string)
	additionalFlags := ctx.Value(cmdAdditionalFlagsCtxKey{}).(map[string]string)
	stdIn := ctx.Value(cmdStdInCtxKey{}).(*bytes.Buffer)

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

	if stdIn.Len() == 0 {
		args = append(args, "--no-input")
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
