package cli_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func iRunExecute(ctx context.Context, recipe string) (context.Context, error) {
	projectDir := ctx.Value(projectDirectoryPathCtxKey{}).(string)
	optionalFlags, flagsAreSet := ctx.Value(cmdOptionalFlagsCtxKey{}).(map[string]string)

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

	if flagsAreSet && optionalFlags != nil {
		for name, value := range optionalFlags {
			args = append(args, fmt.Sprintf("--%s=%s", name, value))
		}
	}

	cmd.SetArgs(args)
	cmd.Execute()
	return ctx, nil
}

func iExecuteRemoteRecipe(ctx context.Context, repository string) (context.Context, error) {
	ociRegistry := ctx.Value(ociRegistryCtxKey{}).(OCIRegistry)
	configDir, configFileExists := ctx.Value(dockerConfigDirectoryPathCtxKey{}).(string)
	optionalFlags, flagsAreSet := ctx.Value(cmdOptionalFlagsCtxKey{}).(map[string]string)
	var flags map[string]string
	if flagsAreSet {
		flags = optionalFlags
	} else {
		flags = make(map[string]string)
	}

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

	ctx = context.WithValue(ctx, cmdOptionalFlagsCtxKey{}, flags)

	return iRunExecute(ctx, url)
}

func executionOfTheRecipeHasSucceeded(ctx context.Context) (context.Context, error) {
	return ctx, expectGivenOutput(ctx, "Recipe executed successfully")
}
