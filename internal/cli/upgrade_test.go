package cli_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/futurice/jalapeno/internal/cli"
	re "github.com/futurice/jalapeno/pkg/recipe"
	"github.com/spf13/pflag"
)

func iUpgradeSauce(ctx context.Context, recipe string) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	projectDir := ctx.Value(projectDirectoryPathCtxKey{}).(string)
	optionalFlagSet, flagsAreSet := ctx.Value(cmdFlagSetCtxKey{}).(*pflag.FlagSet)

	ctx, cmd := wrapCmdOutputs(ctx, cli.NewUpgradeCmd)

	cmd.SetArgs([]string{filepath.Join(recipesDir, recipe)})
	if err := cmd.Flags().Set("dir", projectDir); err != nil {
		return ctx, err
	}

	if flagsAreSet && optionalFlagSet != nil {
		cmd.Flags().AddFlagSet(optionalFlagSet)
	}

	return ctx, cmd.Execute()
}

func conflictsAreReported(ctx context.Context) (context.Context, error) {
	return ctx, expectGivenOutput(ctx, "modified")
}

func noConflictsWereReported(ctx context.Context) (context.Context, error) {
	cmdStdOut := ctx.Value(cmdStdOutCtxKey{}).(*bytes.Buffer)
	cmdStdErr := ctx.Value(cmdStdErrCtxKey{}).(*bytes.Buffer)
	if matched, _ := regexp.MatchString("modified", cmdStdOut.String()); matched {
		return ctx, fmt.Errorf("Conflict in recipe\nstdout:\n%s\n\nstderr:\n%s\n", cmdStdOut, cmdStdErr)
	}
	return ctx, nil
}

func iChangeProjectFileToContain(ctx context.Context, filename, content string) (context.Context, error) {
	projectDir := ctx.Value(projectDirectoryPathCtxKey{}).(string)
	if err := os.WriteFile(filepath.Join(projectDir, filename), []byte(content), 0644); err != nil {
		return ctx, err
	}
	return ctx, nil
}

func iChangeRecipeTemplateToRender(ctx context.Context, recipeName, filename, content string) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	templateFilePath := filepath.Join(recipesDir, recipeName, "templates", filename)
	if err := os.WriteFile(templateFilePath, []byte(content), 0644); err != nil {
		return ctx, err
	}
	return ctx, nil
}

func iChangeRecipeToVersion(ctx context.Context, recipeName, version string) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	recipeFile := filepath.Join(recipesDir, recipeName, re.RecipeFileName+re.YAMLExtension)
	recipeData, err := os.ReadFile(recipeFile)
	if err != nil {
		return ctx, err
	}

	newData := strings.Replace(string(recipeData), "v0.0.1", version, 1)

	if err := os.WriteFile(filepath.Join(recipesDir, recipeName, re.RecipeFileName+re.YAMLExtension), []byte(newData), 0644); err != nil {
		return ctx, err
	}

	return ctx, nil
}
