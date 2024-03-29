package cli_test

import (
	"context"
	"fmt"

	"github.com/cucumber/godog"
	re "github.com/futurice/jalapeno/pkg/recipe"
)

func AddCheckSteps(s *godog.ScenarioContext) {
	s.Step(`^I check new versions for recipes$`, iRunCheck)
	s.Step(`^I check new versions for recipe "([^"]*)"$`, iRunCheckForRecipe)
	s.Step(`^I check and upgrade new versions for recipes$`, iRunCheckWithUpgrade)
	s.Step(`^I check and upgrade new version for recipe "([^"]*)"$`, iRunCheckWithUpgradeForRecipe)
	s.Step(`^I check new versions for recipe "([^"]*)" from the local OCI repository "([^"]*)"$`, iRunCheckForRecipeFrom)
	s.Step(`^the source of the sauce with recipe "([^"]*)" is in the local OCI registry$`, sourceOfTheSauceIsTheLocalOCIRegistry)
}

func iRunCheck(ctx context.Context) (context.Context, error) {
	ociRegistry := ctx.Value(ociRegistryCtxKey{}).(OCIRegistry)
	projectDir := ctx.Value(projectDirectoryPathCtxKey{}).(string)

	args := []string{
		"check",
		fmt.Sprintf("--dir=%s", projectDir),
	}

	if ociRegistry.TLSEnabled {
		args = append(args, "--insecure=true")
	} else {
		args = append(args, "--plain-http=true")
	}

	return executeCLI(ctx, args...)
}

func iRunCheckForRecipe(ctx context.Context, recipe string) (context.Context, error) {
	additionalFlags := ctx.Value(cmdAdditionalFlagsCtxKey{}).(map[string]string)
	additionalFlags["recipe"] = recipe

	return iRunCheck(ctx)
}

func iRunCheckForRecipeFrom(ctx context.Context, recipe, from string) (context.Context, error) {
	ociRegistry := ctx.Value(ociRegistryCtxKey{}).(OCIRegistry)
	additionalFlags := ctx.Value(cmdAdditionalFlagsCtxKey{}).(map[string]string)

	additionalFlags["from"] = fmt.Sprintf("oci://%s/%s", ociRegistry.Resource.GetHostPort("5000/tcp"), from)

	return iRunCheckForRecipe(ctx, recipe)
}

func iRunCheckWithUpgrade(ctx context.Context) (context.Context, error) {
	additionalFlags := ctx.Value(cmdAdditionalFlagsCtxKey{}).(map[string]string)
	additionalFlags["upgrade"] = "true"

	return iRunCheck(ctx)
}

func iRunCheckWithUpgradeForRecipe(ctx context.Context, recipe string) (context.Context, error) {
	additionalFlags := ctx.Value(cmdAdditionalFlagsCtxKey{}).(map[string]string)
	additionalFlags["upgrade"] = "true"

	return iRunCheckForRecipe(ctx, recipe)
}

func sourceOfTheSauceIsTheLocalOCIRegistry(ctx context.Context, recipeName string) (context.Context, error) {
	projectDir := ctx.Value(projectDirectoryPathCtxKey{}).(string)
	ociRegistry := ctx.Value(ociRegistryCtxKey{}).(OCIRegistry)
	sauces, err := re.LoadSauces(projectDir)
	if err != nil {
		return ctx, err
	}

	var sauce *re.Sauce
	for _, s := range sauces {
		if s.Recipe.Name == recipeName {
			sauce = s
			break
		}
	}

	sauce.CheckFrom = fmt.Sprintf("oci://%s/%s", ociRegistry.Resource.GetHostPort("5000/tcp"), sauce.Recipe.Name)
	err = sauce.Save(projectDir)
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}
