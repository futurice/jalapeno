package cli_test

import (
	"context"
	"fmt"
	"path/filepath"
)

func iRunTest(ctx context.Context, recipe string) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	additionalFlags := ctx.Value(cmdAdditionalFlagsCtxKey{}).(map[string]string)

	ctx, cmd := wrapCmdOutputs(ctx)

	args := []string{
		"test",
		filepath.Join(recipesDir, recipe),
	}

	for name, value := range additionalFlags {
		args = append(args, fmt.Sprintf("--%s=%s", name, value))
	}

	cmd.SetArgs(args)
	_ = cmd.Execute()
	return ctx, nil
}

func iCreateRecipeTestUsingCLI(ctx context.Context, recipe string) (context.Context, error) {
	ctx = context.WithValue(
		ctx,
		cmdAdditionalFlagsCtxKey{},
		map[string]string{"create": "true"},
	)

	return iRunTest(ctx, recipe)
}
