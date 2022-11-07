package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/cucumber/godog"
)

type projectDirectoryPathCtxKey struct{}
type recipesDirectoryPathCtxKey struct{}
type templateCtxKey struct{}

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
	outputBasePath = projectDir
	executeFunc(nil, []string{filepath.Join(recipesDir, recipe)})
	return ctx, nil
}

func theProjectDirectoryShouldContainFile(recipe string) error {
	return godog.ErrPending
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

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			s.Step(`^a project directory$`, aProjectDirectory)
			s.Step(`^a recipes directory$`, aRecipesDirectory)
			s.Step(`^a recipe "([^"]*)" that generates file "([^"]*)"$`, aRecipeThatGeneratesFile)
			s.Step(`^I execute recipe "([^"]*)"$`, iExecuteRecipe)
			s.Step(`^the project directory should contain file "([^"]*)"$`, theProjectDirectoryShouldContainFile)
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
