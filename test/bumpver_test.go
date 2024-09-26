package cli_test

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/cucumber/godog"
	"github.com/futurice/jalapeno/pkg/recipe"
)

func AddBumpverSteps(s *godog.ScenarioContext) {
	s.Step(`^I bump recipe "([^"]*)" version to "([^"]*)" with message "([^"]*)"$`, iRunBumpver)
	s.Step(`^recipe "([^"]*)" has version "([^"]*)"`, iCheckVersionNumber)
	s.Step(`^recipe "([^"]*)" has changelog message "([^"]*)"`, iCheckChangelogMsg)
	s.Step(`^recipe "([^"]*)" contains changelog with (\d+) entries`, iVerifyChangelog)
	s.Step(`^first entry in recipe "([^"]*)" changelog has message "([^"]*)"`, iCheckFirstChangelogEntry)
}

func iRunBumpver(ctx context.Context, recipeName, version, msg string) (context.Context, error) {
	recipesDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	re, err := recipe.LoadRecipe(filepath.Join(recipesDir, recipeName))

	if err != nil {
		return ctx, err
	}

	args := []string{
		"bumpver",
		fmt.Sprintf("%s/%s", recipesDir, re.Name),
		fmt.Sprintf("-v=%s", version),
		fmt.Sprintf("-m=%s", msg),
	}

	return executeCLI(ctx, args...)
}

func iCheckVersionNumber(ctx context.Context, recipeName, vers string) (context.Context, error) {
	recipeDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	re, err := recipe.LoadRecipe(filepath.Join(recipeDir, recipeName))

	if err != nil {
		return ctx, err
	}

	if re.Version != vers {
		return ctx, fmt.Errorf("expected version %s, actual: %s", vers, re.Version)
	}

	return ctx, nil
}

func iCheckChangelogMsg(ctx context.Context, recipeName, msg string) (context.Context, error) {
	recipeDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	re, err := recipe.LoadRecipe(filepath.Join(recipeDir, recipeName))

	if err != nil {
		return ctx, err
	}

	if re.Changelog[re.Version] != msg {
		return ctx, fmt.Errorf("expected changelog message %s, actual: %s", msg, re.Changelog[re.Version])
	}

	return ctx, nil
}

func iCheckFirstChangelogEntry(ctx context.Context, recipeName, msg string) (context.Context, error) {
	recipeDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	re, err := recipe.LoadRecipe(filepath.Join(recipeDir, recipeName))

	if err != nil {
		return ctx, err
	}

	if re.Changelog["v0.0.1"] != msg {
		return ctx, fmt.Errorf("expected first changelog message to be %s, actual: %s", msg, re.Changelog["v0.0.1"])
	}

	return ctx, nil
}

func iVerifyChangelog(ctx context.Context, recipeName string, entries int) (context.Context, error) {
	recipeDir := ctx.Value(recipesDirectoryPathCtxKey{}).(string)
	re, err := recipe.LoadRecipe(filepath.Join(recipeDir, recipeName))

	if err != nil {
		return ctx, err
	}

	if len(re.Changelog) != entries {
		return ctx, fmt.Errorf("expected changelog to have %d entries, actual: %d", entries, len(re.Changelog))
	}

	return ctx, nil
}
