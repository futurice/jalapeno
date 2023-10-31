package cli

import (
	"context"

	"github.com/futurice/jalapeno/internal/cli/option"
	"github.com/futurice/jalapeno/pkg/oci"
	"github.com/spf13/cobra"
)

type pushOptions struct {
	RecipePath string
	TargetURL  string
	option.OCIRepository
	option.Common
}

func NewPushCmd() *cobra.Command {
	var opts pushOptions
	var cmd = &cobra.Command{
		Use:   "push RECIPE_PATH TARGET_URL",
		Short: "Push a recipe to OCI repository",
		Long:  "Push a recipe to OCI repository (e.g. Docker registry). You can authenticate by using the --username and --password flags or logging in first with `docker login`.",
		Args:  cobra.ExactArgs(2),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.RecipePath = args[0]
			opts.TargetURL = args[1]
			return option.Parse(&opts)
		},
		Run: func(cmd *cobra.Command, args []string) {
			runPush(cmd, opts)
		},
		Example: `# Push recipe to OCI repository
jalapeno push path/to/recipe ghcr.io/user/recipe:latest

# Push recipe to OCI repository with inline authentication
jalapeno push path/to/recipe oci://ghcr.io/user/my-recipe:latest --username user --password pass

# Push recipe to OCI repository with Docker authentication
docker login ghcr.io
jalapeno push path/to/recipe oci://ghcr.io/user/my-recipe:latest`,
	}

	if err := option.ApplyFlags(&opts, cmd.Flags()); err != nil {
		return nil
	}

	return cmd
}

func runPush(cmd *cobra.Command, opts pushOptions) {
	ctx := context.Background()

	err := oci.PushRecipe(ctx, opts.RecipePath, opts.Repository(opts.TargetURL))

	if err != nil {
		cmd.PrintErrf("Error: %s\n", err)
		return
	}

	cmd.Println("Recipe pushed successfully")
}
