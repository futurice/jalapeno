package cli_test

import (
	"context"
	"fmt"
)

func iRunCreate(ctx context.Context, recipe string) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)

	ctx, cmd := wrapCmdOutputs(ctx)

	args := []string{
		"create",
		recipe,
		fmt.Sprintf("--dir=%s", recipesDir),
	}

	cmd.SetArgs(args)
	cmd.Execute()
	return ctx, nil
}
