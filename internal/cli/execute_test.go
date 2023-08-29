package cli_test

import (
	"context"
	"path/filepath"

	"github.com/futurice/jalapeno/internal/cli"
)

func iRunExecute(ctx context.Context, recipe string) (context.Context, error) {
	projectDir := ctx.Value(projectDirectoryPathCtxKey{}).(string)
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	optionalFlags, flagsAreSet := ctx.Value(cmdOptionalFlagsCtxKey{}).(map[string]string)

	ctx, cmd := wrapCmdOutputs(ctx, cli.NewExecuteCmd)

	cmd.SetArgs([]string{filepath.Join(recipesDir, recipe)})
	if err := cmd.Flags().Set("dir", projectDir); err != nil {
		return ctx, err
	}

	if flagsAreSet && optionalFlags != nil {
		for name, value := range optionalFlags {
			cmd.Flags().Set(name, value)
		}
	}

	return ctx, cmd.Execute()
}

func executionOfTheRecipeHasSucceeded(ctx context.Context) (context.Context, error) {
	return ctx, expectGivenOutput(ctx, "Recipe executed successfully")
}
