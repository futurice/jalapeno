package main

import (
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
	re, err := recipe.LoadRecipe(opts.RecipePath)
	if err != nil {
		cmd.PrintErrf("Can't load the recipe: %v\n", err)
		return
	}

	if len(re.Tests) == 0 {
		cmd.Println("No tests specified")
		return
	}

	errs := re.RunTests()
	errFound := false
	for i, err := range errs {
		if err == nil {
			continue
		}
		cmd.PrintErrf("Test %s failed: %v\n", re.Tests[i].Name, err)
		errFound = true
	}

	if !errFound {
		cmd.Println("Tests passed successfully!")
	}
}
