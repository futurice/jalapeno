package cli_test

import (
	"context"
	"os"
	"path/filepath"

	"github.com/futurice/jalapeno/internal/cli"
)

func iPullRecipe(ctx context.Context, recipeName, repoName string) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	registry := ctx.Value(ociRegistryCtxKey{}).(OCIRegistry)
	configDir, configFileExists := ctx.Value(dockerConfigDirectoryPathCtxKey{}).(string)
	optionalFlags, flagsAreSet := ctx.Value(cmdOptionalFlagsCtxKey{}).(map[string]string)

	ctx, cmd := wrapCmdOutputs(ctx, cli.NewPullCmd)

	cmd.SetArgs([]string{filepath.Join(registry.Resource.GetHostPort("5000/tcp"), repoName)})

	flags := cmd.Flags()
	if err := flags.Set("dir", recipesDir); err != nil {
		return ctx, err
	}
	if registry.TLSEnabled {
		// Allow self-signed certificates
		if err := flags.Set("insecure", "true"); err != nil {
			return ctx, err
		}
	} else {
		if err := flags.Set("plain-http", "true"); err != nil {
			return ctx, err
		}
	}

	if registry.AuthEnabled {
		if err := flags.Set("username", "foo"); err != nil {
			return ctx, err
		}
		if err := flags.Set("password", "bar"); err != nil {
			return ctx, err
		}
	}

	if configFileExists && os.Getenv("DOCKER_CONFIG") == "" {
		if err := flags.Set("registry-config", filepath.Join(configDir, DOCKER_CONFIG_FILENAME)); err != nil {
			return ctx, err
		}
	}

	if flagsAreSet && optionalFlags != nil {
		for name, value := range optionalFlags {
			if err := flags.Set(name, value); err != nil {
				return ctx, err
			}
		}
	}

	cmd.Execute()
	return ctx, nil
}
