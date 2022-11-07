package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/cucumber/godog"
	"github.com/spf13/cobra"
)

type projectDirectoryPathCtxKey struct{}
type recipesDirectoryPathCtxKey struct{}
type recipeStdoutCtxKey struct{}
type recipeStderrCtxKey struct{}

/*
 * UTILITIES
 */

type outputCapturingWriter struct {
	output *string
}

func (o outputCapturingWriter) Write(p []byte) (int, error) {
	*(o.output) = fmt.Sprint(*(o.output), string(p))
	return len(p), nil
}

func newOutputCapturingExecuteCmd(stdout, stderr *string) *cobra.Command {
	w_out := outputCapturingWriter{output: stdout}
	w_err := outputCapturingWriter{output: stderr}
	cmd := newExecuteCmd()
	cmd.SetOut(w_out)
	cmd.SetErr(w_err)
	return cmd
}

/*
 * STEP DEFINITIONS
 */

func aProjectDirectory(ctx context.Context) (context.Context, error) {
	dir, err := os.MkdirTemp("", "jalapeno-test-project")
	if err != nil {
		return ctx, err
	}

	return context.WithValue(ctx, projectDirectoryPathCtxKey{}, dir), nil
}

func aRecipesDirectory(ctx context.Context) (context.Context, error) {
	dir, err := os.MkdirTemp("", "jalapeno-test-recipes")
	if err != nil {
		return ctx, err
	}

	return context.WithValue(ctx, recipesDirectoryPathCtxKey{}, dir), nil
}

func aRecipeThatGeneratesFile(ctx context.Context, recipe, filename string) (context.Context, error) {
	dir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	if err := os.MkdirAll(filepath.Join(dir, recipe, "templates"), 0755); err != nil {
		return ctx, err
	}
	template := "name: %[1]s\nversion: v0.0.1\ndescription: %[1]s"
	if err := os.WriteFile(filepath.Join(dir, recipe, "recipe.yml"), []byte(fmt.Sprintf(template, recipe)), 0660); err != nil {
		return ctx, err
	}
	if err := os.WriteFile(filepath.Join(dir, recipe, "templates", filename), []byte(recipe), 0660); err != nil {
		return ctx, err
	}
	return context.WithValue(ctx, recipesDirectoryPathCtxKey{}, dir), nil
}

func iExecuteRecipe(ctx context.Context, recipe string) (context.Context, error) {
	projectDir := ctx.Value(projectDirectoryPathCtxKey{}).(string)
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	recipeStdout := ""
	recipeStderr := ""
	cmd := newOutputCapturingExecuteCmd(&recipeStdout, &recipeStderr)
	cmd.Flags().Set("output", projectDir)
	executeFunc(cmd, []string{filepath.Join(recipesDir, recipe)})
	return context.WithValue(
		context.WithValue(ctx, recipeStdoutCtxKey{}, recipeStdout),
		recipeStderrCtxKey{},
		recipeStderr,
	), nil
}

func theProjectDirectoryShouldContainFile(ctx context.Context, filename string) error {
	dir := ctx.Value(projectDirectoryPathCtxKey{}).(string)
	info, err := os.Stat(filepath.Join(dir, filename))
	if err == nil && !info.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", filename)
	}
	return err
}

func cleanTempDirs(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
	if dir := ctx.Value(projectDirectoryPathCtxKey{}); dir != nil {
		fmt.Printf("Removing %s\n", dir)
		os.RemoveAll(dir.(string))
	}
	if dir := ctx.Value(recipesDirectoryPathCtxKey{}); dir != nil {
		fmt.Printf("Removing %s\n", dir)
		os.RemoveAll(dir.(string))
	}
	return ctx, err
}

func executionOfRecipeHasFailedWithError(ctx context.Context, _recipe, errorMessage string) (context.Context, error) {
	out := ctx.Value(recipeStderrCtxKey{}).(string)
	if match, err := regexp.Match(errorMessage, []byte(out)); !match {
		if err != nil {
			return ctx, err
		}
		return ctx, fmt.Errorf("'%s' not found in stderr", errorMessage)
	}
	return ctx, nil
}

func iChangeRecipeToVersion(arg1, arg2 string) error {
	return godog.ErrPending
}

func iUpgradeRecipe(arg1 string) error {
	return godog.ErrPending
}

func theProjectDirectoryShouldContainFileWith(arg1, arg2 string) error {
	return godog.ErrPending
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			s.Step(`^a project directory$`, aProjectDirectory)
			s.Step(`^a recipes directory$`, aRecipesDirectory)
			s.Step(`^a recipe "([^"]*)" that generates file "([^"]*)"$`, aRecipeThatGeneratesFile)
			s.Step(`^I execute recipe "([^"]*)"$`, iExecuteRecipe)
			s.Step(`^the project directory should contain file "([^"]*)"$`, theProjectDirectoryShouldContainFile)
			s.Step(`^the project directory should contain file "([^"]*)" with "([^"]*)"$`, theProjectDirectoryShouldContainFileWith)
			s.Step(`^execution of recipe "([^"]*)" has failed with error "([^"]*)"$`, executionOfRecipeHasFailedWithError)
			s.Step(`^I change recipe "([^"]*)" to version "([^"]*)"$`, iChangeRecipeToVersion)
			s.Step(`^I upgrade recipe "([^"]*)"$`, iUpgradeRecipe)
			s.After(cleanTempDirs)
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
