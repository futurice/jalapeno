package cli_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/futurice/jalapeno/internal/cli"
)

func iRunExecute(ctx context.Context, recipe string) (context.Context, error) {
	projectDir := ctx.Value(projectDirectoryPathCtxKey{}).(string)
	optionalFlags, flagsAreSet := ctx.Value(cmdOptionalFlagsCtxKey{}).(map[string]string)

	ctx, cmd := wrapCmdOutputs(ctx, cli.NewExecuteCmd)

	var url string
	if strings.HasPrefix(recipe, "oci://") {
		url = recipe
	} else {
		recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
		url = filepath.Join(recipesDir, recipe)
	}

	cmd.SetArgs([]string{url})

	flags := cmd.Flags()
	if err := flags.Set("dir", projectDir); err != nil {
		return ctx, err
	}

	if flagsAreSet && optionalFlags != nil {
		for name, value := range optionalFlags {
			if err := flags.Set(name, value); err != nil {
				return ctx, err
			}
		}
	}

	return ctx, cmd.Execute()
}

func iExecuteRemoteRecipe(ctx context.Context, repository string) (context.Context, error) {
	registry := ctx.Value(ociRegistryCtxKey{}).(OCIRegistry)
	configDir, configFileExists := ctx.Value(dockerConfigDirectoryPathCtxKey{}).(string)
	optionalFlags, flagsAreSet := ctx.Value(cmdOptionalFlagsCtxKey{}).(map[string]string)
	var flags map[string]string
	if flagsAreSet {
		flags = optionalFlags
	} else {
		flags = make(map[string]string)
	}

	x := fmt.Sprintf("oci://%s/%s", registry.Resource.GetHostPort("5000/tcp"), repository)

	if registry.TLSEnabled {
		// Allow self-signed certificates
		flags["insecure"] = "true"
	} else {
		flags["plain-http"] = "true"
	}

	if registry.AuthEnabled {
		flags["username"] = "foo"
		flags["password"] = "bar"
	}

	if configFileExists && os.Getenv("DOCKER_CONFIG") == "" {
		flags["registry-config"] = filepath.Join(configDir, DOCKER_CONFIG_FILENAME)
	}

	ctx = context.WithValue(ctx, cmdOptionalFlagsCtxKey{}, flags)

	return iRunExecute(ctx, x)
}

func executionOfTheRecipeHasSucceeded(ctx context.Context) (context.Context, error) {
	return ctx, expectGivenOutput(ctx, "Recipe executed successfully")
}
