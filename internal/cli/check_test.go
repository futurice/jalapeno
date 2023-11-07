package cli_test

import (
	"context"
	"fmt"

	re "github.com/futurice/jalapeno/pkg/recipe"
)

func iRunCheck(ctx context.Context, recipe string) (context.Context, error) {
	projectDir := ctx.Value(projectDirectoryPathCtxKey{}).(string)
	ociRegistry := ctx.Value(ociRegistryCtxKey{}).(OCIRegistry)
	optionalFlags, flagsAreSet := ctx.Value(cmdOptionalFlagsCtxKey{}).(map[string]string)

	ctx, cmd := wrapCmdOutputs(ctx)

	args := []string{
		"check",
		fmt.Sprintf("--recipe=%s", recipe),
		fmt.Sprintf("--dir=%s", projectDir),
	}

	if ociRegistry.TLSEnabled {
		args = append(args, "--insecure=true")
	} else {
		args = append(args, "--plain-http=true")
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

func newRecipeVersionsWereFound(ctx context.Context) (context.Context, error) {
	return ctx, expectGivenOutput(ctx, "New versions found")
}

func noNewRecipeVersionsWereFound(ctx context.Context) (context.Context, error) {
	return ctx, expectGivenOutput(ctx, "No new versions found")
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
