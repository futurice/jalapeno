package cli_test

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/futurice/jalapeno/internal/cli"
	re "github.com/futurice/jalapeno/pkg/recipe"
)

func iRunCheck(ctx context.Context, recipe string) (context.Context, error) {
	projectDir := ctx.Value(projectDirectoryPathCtxKey{}).(string)
	ociRegistry := ctx.Value(ociRegistryCtxKey{}).(OCIRegistry)
	optionalFlags, flagsAreSet := ctx.Value(cmdOptionalFlagsCtxKey{}).(map[string]string)

	ctx, cmd := wrapCmdOutputs(ctx, cli.NewCheckCmd)

	cmd.SetArgs([]string{recipe})

	flags := cmd.Flags()
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

	return ctx, cmd.Execute()
}

func newRecipeVersionsWereFound(ctx context.Context) (context.Context, error) {
	return ctx, expectGivenOutput(ctx, "New versions found")
}

func noNewRecipeVersionsWereFound(ctx context.Context) (context.Context, error) {
	return ctx, expectGivenOutput(ctx, "No new versions found")
}

func sourceOfTheRecipeIsTheLocalOCIRegistry(ctx context.Context, recipeName string) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	ociRegistry := ctx.Value(ociRegistryCtxKey{}).(OCIRegistry)
	recipePath := filepath.Join(recipesDir, recipeName)
	recipe, err := re.LoadRecipe(recipePath)
	if err != nil {
		return ctx, err
	}

	url := fmt.Sprintf("%s/%s", ociRegistry.Resource.GetHostPort("5000/tcp"), recipe.Name)
	recipe.Sources = append(recipe.Sources, url)
	err = recipe.Save(recipesDir)
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}