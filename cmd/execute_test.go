package main

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
)

func executionOfTheRecipeHasSucceeded(ctx context.Context) (context.Context, error) {
	cmdStdOut := ctx.Value(cmdStdOutCtxKey{}).(string)
	cmdStdErr := ctx.Value(cmdStdErrCtxKey{}).(string)
	if matched, _ := regexp.MatchString("Recipe executed successfully", cmdStdOut); !matched {
		return ctx, fmt.Errorf("Recipe failed to execute!\nstdout:\n%s\n\nstderr:\n%s\n", cmdStdOut, cmdStdErr)
	}
	return ctx, nil
}

func executionOfTheRecipeHasFailedWithError(ctx context.Context, errorMessage string) (context.Context, error) {
	cmdStdOut := ctx.Value(cmdStdOutCtxKey{}).(string)
	cmdStdErr := ctx.Value(cmdStdErrCtxKey{}).(string)
	if matched, _ := regexp.MatchString(errorMessage, cmdStdErr); !matched {
		return ctx, fmt.Errorf("'%s' not found in stderr.\nstdout:\n%s\n\nstderr:\n%s\n", errorMessage, cmdStdOut, cmdStdErr)
	}
	return ctx, nil
}

func iExecuteRecipe(ctx context.Context, recipe string) (context.Context, error) {
	projectDir := ctx.Value(projectDirectoryPathCtxKey{}).(string)
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)

	cmd, cmdStdOut, cmdStdErr := WrapCmdOutputs(newExecuteCmd)

	cmd.SetArgs([]string{filepath.Join(recipesDir, recipe)})
	if err := cmd.Flags().Set("output", projectDir); err != nil {
		return ctx, err
	}

	if err := cmd.Execute(); err != nil {
		return ctx, err
	}

	ctx = context.WithValue(ctx, cmdStdOutCtxKey{}, cmdStdOut.String())
	ctx = context.WithValue(ctx, cmdStdErrCtxKey{}, cmdStdErr.String())

	return ctx, nil
}
