package main

import (
	"errors"

	"github.com/futurice/jalapeno/cmd/internal/option"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/spf13/cobra"
)

type testOptions struct {
	RecipePath string
	option.Common
}

func newTestCmd() *cobra.Command {
	var opts testOptions
	var cmd = &cobra.Command{
		Use:   "test RECIPE_PATH",
		Short: "Test a recipe",
		Long:  "",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.RecipePath = args[0]
			return option.Parse(&opts)
		},
		Run: func(cmd *cobra.Command, args []string) {
			runTest(cmd, opts)
		},
	}

	if err := option.ApplyFlags(&opts, cmd.Flags()); err != nil {
		return nil
	}

	return cmd
}

func runTest(cmd *cobra.Command, opts testOptions) {
	re, err := recipe.Load(opts.RecipePath)
	if err != nil {
		cmd.PrintErrf("Can't load the recipe: %v\n", err)
		return
	}

	err = re.RunTests()
	if err != nil {
		if errors.Is(err, recipe.ErrNoTestsSpecified) {
			cmd.Println("No tests specified")
			return
		}

		cmd.PrintErrf("Tests failed: %v\n", err)
		return
	}

	cmd.Println("Tests passed successfully!")
}
