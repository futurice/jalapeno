package cli

import (
	"fmt"
	"path/filepath"

	"github.com/futurice/jalapeno/internal/cli/option"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/ui/colors"
	"github.com/spf13/cobra"
)

type validateOptions struct {
	RecipePath string

	option.Common
	option.WorkingDirectory
}

func NewValidateCmd() *cobra.Command {
	var opts validateOptions
	var cmd = &cobra.Command{
		Use:   "validate RECIPE_PATH",
		Short: "Validate a recipe",
		Long:  "Validate a recipe in a local path.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.RecipePath = args[0]
			return option.Parse(&opts)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runValidate(cmd, opts)
			return errorHandler(cmd, err)
		},
		Args:    cobra.ExactArgs(1),
		Example: `jalapeno validate path/to/recipe`,
	}

	if err := option.ApplyFlags(&opts, cmd.Flags()); err != nil {
		return nil
	}

	return cmd
}

func runValidate(cmd *cobra.Command, opts validateOptions) error {
	path := opts.RecipePath
	if !filepath.IsAbs(opts.RecipePath) {
		path = filepath.Join(opts.Dir, opts.RecipePath)
	}

	r, err := recipe.LoadRecipe(path)
	if err != nil {
		return fmt.Errorf("could not load the recipe: %w", err)
	}

	err = r.Validate()
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	cmd.Printf(
		"%s The recipe is valid.\n",
		colors.Green.Render("Success!"),
	)

	return nil
}
