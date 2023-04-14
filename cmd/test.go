package main

import (
	"path/filepath"

	"github.com/futurice/jalapeno/cmd/internal/option"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/spf13/cobra"
)

type testOptions struct {
	RecipePath      string
	UpdateSnapshots bool
	option.Common
}

func newTestCmd() *cobra.Command {
	var opts testOptions
	var cmd = &cobra.Command{
		Use:   "test RECIPE",
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

	cmd.Flags().BoolVarP(&opts.UpdateSnapshots, "update-snapshots", "u", false, "Update file snapshots")

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

	if opts.UpdateSnapshots {
		for i := range re.Tests {
			test := &re.Tests[i]
			sauce, err := re.Execute(test.Values, recipe.TestAnchor)
			if err != nil {
				cmd.PrintErrf("Error: failed to render templates: %s", err)
				return
			}

			test.Files = make(map[string]recipe.TestFile)
			for filename, file := range sauce.Files {
				test.Files[filename] = file.Content
			}
		}

		err := re.Save(filepath.Dir(opts.RecipePath))
		if err != nil {
			cmd.PrintErrf("Error: failed to save recipe: %s", err)
			return
		}

		// TODO: Show which tests were modified
		cmd.Println("Recipe tests updated successfully!")
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
		// TODO: Show pass for each test
		cmd.Println("Tests passed successfully!")
	}
}
