package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/futurice/jalapeno/internal/cli/option"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/spf13/cobra"
)

type pullOptions struct {
	TargetRef string

	option.Common
	option.OCIRepository
	option.WorkingDirectory
}

func NewPullCmd() *cobra.Command {
	var opts pullOptions
	var cmd = &cobra.Command{
		Use:   "pull URL",
		Short: "Pull a recipe from OCI repository",
		Long:  "Pull a recipe from OCI repository and save it locally. You can authenticate by using the `--username` and `--password` flags or logging in first with `docker login`.",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.TargetRef = args[0]
			return option.Parse(&opts)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPull(cmd, opts)
		},
		Example: `# Pull recipe from OCI repository
jalapeno pull ghcr.io/user/recipe:latest

# Pull recipe from OCI repository with inline authentication
jalapeno pull oci://ghcr.io/user/my-recipe:latest --username user --password pass

# Pull recipe from OCI repository with Docker authentication
docker login ghcr.io
jalapeno pull oci://ghcr.io/user/my-recipe:latest

# Pull recipe to different directory
jalapeno pull oci://ghcr.io/user/my-recipe:latest --dir other/dir`,
	}

	if err := option.ApplyFlags(&opts, cmd.Flags()); err != nil {
		return nil
	}

	return cmd
}

func runPull(cmd *cobra.Command, opts pullOptions) error {
	ctx := context.Background()

	err := recipe.SaveRemoteRecipe(ctx, opts.Dir, opts.Repository(opts.TargetRef))

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return fmt.Errorf("recipe not found from %s", opts.TargetRef)
		}

		return err
	}

	cmd.Printf("Recipe pulled %s\n", opts.Colors.Green.Render("successfully!"))
	return nil
}
