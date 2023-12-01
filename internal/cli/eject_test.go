package cli_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/futurice/jalapeno/pkg/recipe"
)

func iRunEject(ctx context.Context) (context.Context, error) {
	projectDir := ctx.Value(projectDirectoryPathCtxKey{}).(string)

	ctx, cmd := wrapCmdOutputs(ctx)
	args := []string{
		"eject",
		fmt.Sprintf("--dir=%s", projectDir),
	}

	cmd.SetArgs(args)
	_ = cmd.Execute()
	return ctx, nil
}

func thereShouldNotBeASauceDirectoryInTheProjectDirectory(ctx context.Context) (context.Context, error) {
	projectDir := ctx.Value(projectDirectoryPathCtxKey{}).(string)
	_, err := os.Stat(filepath.Join(projectDir, recipe.SauceDirName))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ctx, nil
		}
		return ctx, fmt.Errorf("Expected ErrNotExist, got %w", err)
	}
	return ctx, fmt.Errorf("Expected sauce dir not to exist")
}
