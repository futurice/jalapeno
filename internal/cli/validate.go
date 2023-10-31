package cli

import (
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
		Run: func(cmd *cobra.Command, args []string) {
			runValidate(cmd, opts)
		},
		Args:    cobra.ExactArgs(1),
		Example: `jalapeno validate path/to/recipe`,
	}

	if err := option.ApplyFlags(&opts, cmd.Flags()); err != nil {
		return nil
	}

	return cmd
}

func runValidate(cmd *cobra.Command, opts validateOptions) {
	r, err := recipe.LoadRecipe(opts.RecipePath)
	if err != nil {
		cmd.PrintErrf("Error: could not load the recipe: %s\n", err)
	}

	err = r.Validate()
	if err != nil {
		cmd.PrintErrf("Error: validation failed: %s\n", err)
	}

	cmd.Println("Validation ok")
}
