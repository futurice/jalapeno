package cli

import (
	"context"
	"fmt"

	"github.com/futurice/jalapeno/internal/cli/option"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/ui/colors"
	"github.com/spf13/cobra"
)

type pushOptions struct {
	RecipePath   string
	TargetURL    string
	PushToLatest bool

	option.Common
	option.OCIRepository
	option.Timeout
}

func NewPushCmd() *cobra.Command {
	var opts pushOptions
	var cmd = &cobra.Command{
		Use:   "push RECIPE_PATH TARGET_URL",
		Short: "Push a recipe to OCI repository",
		Long:  "Push a recipe to OCI repository (e.g. Docker registry). The version of the recipe will be used as a tag for the image. You can authenticate by using the `--username` and `--password` flags or logging in first with `docker login`.",
		Args:  cobra.ExactArgs(2),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.RecipePath = args[0]
			opts.TargetURL = args[1]
			return option.Parse(&opts)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runPush(cmd, opts)
			return errorHandler(cmd, err)
		},
		Example: `# Push recipe to OCI repository
jalapeno push path/to/recipe ghcr.io/user/recipe

# Push recipe to OCI repository with inline authentication
jalapeno push path/to/recipe oci://ghcr.io/user/my-recipe --username user --password pass

# Push recipe to OCI repository with Docker authentication
docker login ghcr.io
jalapeno push path/to/recipe oci://ghcr.io/user/my-recipe`,
	}

	cmd.Flags().BoolVarP(&opts.PushToLatest, "latest", "l", false, "Additionally push the recipe to 'latest' tag")

	if err := option.ApplyFlags(&opts, cmd.Flags()); err != nil {
		return nil
	}

	return cmd
}

func runPush(cmd *cobra.Command, opts pushOptions) error {
	ctx, cancel := context.WithTimeout(cmd.Context(), opts.Timeout.Duration) // nolint:staticcheck
	defer cancel()

	err := recipe.PushRecipe(ctx, opts.RecipePath, opts.Repository(opts.TargetURL), opts.PushToLatest)

	if err != nil {
		return fmt.Errorf("failed to push recipe: %w", err)
	}

	cmd.Printf("Recipe pushed %s\n", colors.Green.Render("successfully!"))
	return nil
}
