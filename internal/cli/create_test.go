package cli_test

import (
	"context"

	"github.com/futurice/jalapeno/internal/cli"
)

func iRunCreate(ctx context.Context, recipe string) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)

	ctx, cmd := wrapCmdOutputs(ctx, cli.NewCreateCmd)

	cmd.SetArgs([]string{recipe})

	flags := cmd.Flags()
	if err := flags.Set("dir", recipesDir); err != nil {
		return ctx, err
	}

	return ctx, cmd.Execute()
}
