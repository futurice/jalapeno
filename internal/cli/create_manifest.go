package cli

import (
	"fmt"
	"path/filepath"

	"github.com/futurice/jalapeno/internal/cli/option"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/recipeutil"
	"github.com/spf13/cobra"
)

type createManifestOptions struct {
	option.Common
	option.WorkingDirectory
}

func NewCreateManifestCmd() *cobra.Command {
	var opts createManifestOptions
	var cmd = &cobra.Command{
		Use:     "manifest",
		Short:   "Create a manifest",
		Example: `jalapeno create manifest`,
		Args:    cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return option.Parse(&opts)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runCreateManifest(cmd, opts)
			return errorHandler(cmd, err)
		},
	}

	if err := option.ApplyFlags(&opts, cmd.Flags()); err != nil {
		return nil
	}

	return cmd
}

func runCreateManifest(cmd *cobra.Command, opts createManifestOptions) error {
	m := recipeutil.CreateExampleManifest()

	err := m.Save(filepath.Join(opts.Dir, fmt.Sprintf(recipe.ManifestFileName+recipe.YAMLExtension)))
	if err != nil {
		return fmt.Errorf("can not save recipe to the directory: %w", err)
	}

	cmd.Printf(
		"Manifest created %s\n",
		ColorGreen.Render("successfully!"),
	)

	return nil
}
