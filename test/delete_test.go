package cli_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cucumber/godog"
	"github.com/futurice/jalapeno/pkg/recipe"
)

func AddDeleteSteps(s *godog.ScenarioContext) {
	s.Step(`^I delete the sauce from the index (\d)$`, iRunDelete)
	s.Step(`^I delete all sauces from the project$`, iRunDeleteAll)
	s.Step(`^there should not be a sauce directory in the project directory$`, thereShouldNotBeASauceDirInTheProjectDir)
}

func iRunDelete(ctx context.Context, i int) (context.Context, error) {
	projectDir := ctx.Value(projectDirectoryPathCtxKey{}).(string)
	sauces, err := recipe.LoadSauces(projectDir)
	if err != nil {
		return ctx, err
	}

	if i < 0 || i >= len(sauces) {
		return ctx, fmt.Errorf("invalid sauce index: %d", i)
	}

	return executeCLI(ctx,
		"delete",
		sauces[i].ID.String(),
		fmt.Sprintf("--dir=%s", projectDir),
	)
}

func iRunDeleteAll(ctx context.Context) (context.Context, error) {
	projectDir := ctx.Value(projectDirectoryPathCtxKey{}).(string)

	return executeCLI(ctx,
		"delete",
		"--all",
		fmt.Sprintf("--dir=%s", projectDir),
	)
}

func thereShouldNotBeASauceDirInTheProjectDir(ctx context.Context) (context.Context, error) {
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
