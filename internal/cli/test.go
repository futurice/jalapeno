package cli

import (
	"path/filepath"

	"github.com/futurice/jalapeno/internal/cli/option"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/recipeutil"
	"github.com/spf13/cobra"
)

type testOptions struct {
	RecipePath      string
	UpdateSnapshots bool
	Create          bool
	option.Common
}

func NewTestCmd() *cobra.Command {
	var opts testOptions
	var cmd = &cobra.Command{
		Use:   "test RECIPE_PATH",
		Short: "Run tests for the recipe",
		Long:  "TODO",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.RecipePath = args[0]
			return option.Parse(&opts)
		},
		Run: func(cmd *cobra.Command, args []string) {
			runTest(cmd, opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.UpdateSnapshots, "update-snapshots", "u", false, "Update test file snapshots")
	cmd.Flags().BoolVarP(&opts.Create, "create", "c", false, "Create a new test case")

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

	if opts.Create {
		test := *recipeutil.CreateExampleTest()

		if len(re.Tests) > 0 {
			re.Tests = append(re.Tests, test)
		} else {
			re.Tests = []recipe.Test{test}
		}

		err := re.Save(filepath.Dir(opts.RecipePath))
		if err != nil {
			cmd.PrintErrf("Error: failed to save recipe: %s", err)
			return
		}

		cmd.Println("Test created")
		return
	}

	if len(re.Tests) == 0 {
		cmd.Println("No tests specified")
		return
	}

	if opts.UpdateSnapshots {
		for i := range re.Tests {
			test := &re.Tests[i]
			sauce, err := re.Execute(test.Values, recipe.TestID)
			if err != nil {
				cmd.PrintErrf("Error: failed to render templates: %s", err)
				return
			}

			test.Files = make(map[string][]byte)
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
		cmd.Println("Recipe tests updated successfully")
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
		cmd.Println("Tests passed successfully")
	}
}
