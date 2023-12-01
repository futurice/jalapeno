package cli_test

import (
	"context"
	"fmt"

	re "github.com/futurice/jalapeno/pkg/recipe"
)

func iRunCheck(ctx context.Context) (context.Context, error) {
	projectDir := ctx.Value(projectDirectoryPathCtxKey{}).(string)
	ociRegistry := ctx.Value(ociRegistryCtxKey{}).(OCIRegistry)
	additionalFlags := ctx.Value(cmdAdditionalFlagsCtxKey{}).(map[string]string)

	ctx, cmd := wrapCmdOutputs(ctx)

	args := []string{
		"check",
		fmt.Sprintf("--dir=%s", projectDir),
	}

	if ociRegistry.TLSEnabled {
		args = append(args, "--insecure=true")
	} else {
		args = append(args, "--plain-http=true")
	}

	for name, value := range additionalFlags {
		args = append(args, fmt.Sprintf("--%s=%s", name, value))
	}

	cmd.SetArgs(args)
	_ = cmd.Execute()
	return ctx, nil
}

func iRunCheckForRecipe(ctx context.Context, recipe string) (context.Context, error) {
	additionalFlags := ctx.Value(cmdAdditionalFlagsCtxKey{}).(map[string]string)
	additionalFlags["recipe"] = recipe

	ctx = context.WithValue(ctx, cmdAdditionalFlagsCtxKey{}, additionalFlags)
	return iRunCheck(ctx)
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
