package cli_test

import (
	"context"
	"path/filepath"

	"github.com/futurice/jalapeno/internal/cli"
	"github.com/spf13/pflag"
)

func iPushRecipe(ctx context.Context, recipeName string) (context.Context, error) {
	return pushRecipe(ctx, recipeName)
}

func pushRecipe(ctx context.Context, recipeName string) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	ociRegistry := ctx.Value(ociRegistryCtxKey{}).(OCIRegistry)
	configDir, configFileExists := ctx.Value(dockerConfigDirectoryPathCtxKey{}).(string)
	optionalFlagSet, flagsAreSet := ctx.Value(cmdFlagSetCtxKey{}).(*pflag.FlagSet)

	ctx, cmd := wrapCmdOutputs(ctx, cli.NewPushCmd)

	cmd.SetArgs([]string{filepath.Join(recipesDir, recipeName), filepath.Join(ociRegistry.Resource.GetHostPort("5000/tcp"), recipeName)})

	flags := cmd.Flags()
	if ociRegistry.TLSEnabled {
		if err := flags.Set("insecure", "true"); err != nil {
			return ctx, err
		}
	} else {
		if err := flags.Set("plain-http", "true"); err != nil {
			return ctx, err
		}
	}

	if ociRegistry.AuthEnabled {
		if err := flags.Set("username", "foo"); err != nil {
			return ctx, err
		}
		if err := flags.Set("password", "bar"); err != nil {
			return ctx, err
		}
	}

	if configFileExists {
		if err := flags.Set("registry-config", filepath.Join(configDir, DOCKER_CONFIG_FILENAME)); err != nil {
			return ctx, err
		}
	}

	if flagsAreSet && optionalFlagSet != nil {
		cmd.Flags().AddFlagSet(optionalFlagSet)
	}

	return ctx, cmd.Execute()
}

func pushOfTheRecipeWasSuccessful(ctx context.Context) (context.Context, error) {
	return ctx, expectGivenOutput(ctx, "") // TODO
}
