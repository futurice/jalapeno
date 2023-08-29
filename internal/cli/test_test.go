package cli_test

import (
	"context"
	"path/filepath"

	"github.com/futurice/jalapeno/internal/cli"
)

func iRunTest(ctx context.Context, recipe string) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	optionalFlags, flagsAreSet := ctx.Value(cmdOptionalFlagsCtxKey{}).(map[string]string)

	ctx, cmd := wrapCmdOutputs(ctx, cli.NewTestCmd)

	cmd.SetArgs([]string{filepath.Join(recipesDir, recipe)})

	if flagsAreSet && optionalFlags != nil {
		for name, value := range optionalFlags {
			cmd.Flags().Set(name, value)
		}
	}

	return ctx, cmd.Execute()
}

func iCreateRecipeTestUsingCLI(ctx context.Context, recipe string) (context.Context, error) {
	ctx = context.WithValue(
		ctx,
		cmdOptionalFlagsCtxKey{},
		map[string]string{"create": "true"},
	)

	return iRunTest(ctx, recipe)
}
