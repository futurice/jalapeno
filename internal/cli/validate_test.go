package cli_test

import (
	"context"
	"fmt"
)

func iRunValidate(ctx context.Context, recipe string) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	optionalFlags, flagsAreSet := ctx.Value(cmdOptionalFlagsCtxKey{}).(map[string]string)

	ctx, cmd := wrapCmdOutputs(ctx)

	args := []string{
		"validate",
		recipe,
		fmt.Sprintf("--dir=%s", recipesDir),
	}

	if flagsAreSet && optionalFlags != nil {
		for name, value := range optionalFlags {
			args = append(args, fmt.Sprintf("--%s=%s", name, value))
		}
	}

	cmd.SetArgs(args)
	_ = cmd.Execute()
	return ctx, nil
}
