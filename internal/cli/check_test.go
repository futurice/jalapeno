package cli_test

import (
	"context"
	"fmt"

	"github.com/futurice/jalapeno/internal/cli"
	re "github.com/futurice/jalapeno/pkg/recipe"
)

func iRunCheck(ctx context.Context, recipe string) (context.Context, error) {
	projectDir := ctx.Value(projectDirectoryPathCtxKey{}).(string)
	ociRegistry := ctx.Value(ociRegistryCtxKey{}).(OCIRegistry)
	optionalFlags, flagsAreSet := ctx.Value(cmdOptionalFlagsCtxKey{}).(map[string]string)

	ctx, cmd := wrapCmdOutputs(ctx, cli.NewCheckCmd)

	flags := cmd.Flags()
	if err := flags.Set("recipe", recipe); err != nil {
		return ctx, err
	}
	if err := flags.Set("dir", projectDir); err != nil {
		return ctx, err
	}
	if ociRegistry.TLSEnabled {
		if err := flags.Set("insecure", "true"); err != nil {
			return ctx, err
		}
	} else {
		if err := flags.Set("plain-http", "true"); err != nil {
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

func newRecipeVersionsWereFound(ctx context.Context) (context.Context, error) {
	return ctx, expectGivenOutput(ctx, "New versions found")
}

func noNewRecipeVersionsWereFound(ctx context.Context) (context.Context, error) {
	return ctx, expectGivenOutput(ctx, "No new versions found")
}

func sourceOfTheSauceIsTheLocalOCIRegistry(ctx context.Context, recipeName string) (context.Context, error) {
	projectDir := ctx.Value(projectDirectoryPathCtxKey{}).(string)
	ociRegistry := ctx.Value(ociRegistryCtxKey{}).(OCIRegistry)
	sauces, err := re.LoadSauces(projectDir)
	if err != nil {
		return ctx, err
	}

	var sauce *re.Sauce
	for _, s := range sauces {
		if s.Recipe.Name == recipeName {
			sauce = s
			break
		}
	}

	sauce.CheckFrom = fmt.Sprintf("oci://%s/%s", ociRegistry.Resource.GetHostPort("5000/tcp"), sauce.Recipe.Name)
	err = sauce.Save(projectDir)
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}
