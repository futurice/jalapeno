package cli

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	re "github.com/futurice/jalapeno/pkg/recipe"
	"github.com/spf13/pflag"
)

func iCheckRecipe(ctx context.Context, recipe string) (context.Context, error) {
	projectDir := ctx.Value(projectDirectoryPathCtxKey{}).(string)
	ociRegistry := ctx.Value(ociRegistryCtxKey{}).(OCIRegistry)
	optionalFlagSet, flagsAreSet := ctx.Value(cmdFlagSetCtxKey{}).(*pflag.FlagSet)

	cmd, cmdStdOut, cmdStdErr := wrapCmdOutputs(newCheckCmd)

	cmd.SetArgs([]string{projectDir, recipe})

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

	if flagsAreSet && optionalFlagSet != nil {
		cmd.Flags().AddFlagSet(optionalFlagSet)
	}

	if err := cmd.Execute(); err != nil {
		return ctx, err
	}

	ctx = context.WithValue(ctx, cmdStdOutCtxKey{}, cmdStdOut.String())
	ctx = context.WithValue(ctx, cmdStdErrCtxKey{}, cmdStdErr.String())

	return ctx, nil
}

func newRecipeVersionsWereFound(ctx context.Context) (context.Context, error) {
	cmdStdOut := ctx.Value(cmdStdOutCtxKey{}).(string)
	cmdStdErr := ctx.Value(cmdStdErrCtxKey{}).(string)

	if cmdStdErr != "" {
		return ctx, fmt.Errorf("command caused unexpected error: %s", cmdStdErr)
	}

	if cmdStdOut == "" {
		return ctx, fmt.Errorf("command output was empty")
	}

	expectedOutput := "New versions found"
	if !strings.Contains(cmdStdOut, expectedOutput) {
		return ctx, fmt.Errorf("command produced unexpected output: Expected: '%s', Actual: '%s'", expectedOutput, cmdStdOut)
	}

	return ctx, nil
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
