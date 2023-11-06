package cli_test

import (
	"context"
	"fmt"

	"github.com/futurice/jalapeno/internal/cli"
)

func iRunCreate(ctx context.Context, recipe string) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)

	ctx, cmd := wrapCmdOutputs(ctx, cli.NewCreateCmd)

	args := []string{
		recipe,
		fmt.Sprintf("--dir=%s", recipesDir),
	}

	cmd.SetArgs(args)
	cmd.Execute()
	return ctx, nil
}
