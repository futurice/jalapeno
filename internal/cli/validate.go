package cli

import (
	"fmt"

	"github.com/futurice/jalapeno/internal/cli/option"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/spf13/cobra"
)

type validateOptions struct {
	RecipePath string
	option.Common
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
			return runValidate(cmd, opts)
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
	r, err := recipe.LoadRecipe(opts.RecipePath)
	if err != nil {
		return fmt.Errorf("could not load the recipe: %w", err)
	}

	err = r.Validate()
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	cmd.Println("Validation ok")
	return nil
}
