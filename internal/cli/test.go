package cli

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/futurice/jalapeno/internal/cli/option"
	"github.com/futurice/jalapeno/pkg/engine"
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
		Long:  "Run tests for the recipe.",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.RecipePath = args[0]
			return option.Parse(&opts)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTest(cmd, opts)
		},
		Example: `# Run recipe tests
jalapeno test path/to/recipe

# Bootstrap a new test case
jalapeno test path/to/recipe --create

# Update test file snapshots with the current outputs
jalapeno test path/to/recipe --update-snapshots`,
	}

	cmd.Flags().BoolVarP(&opts.UpdateSnapshots, "update-snapshots", "u", false, "Update test file snapshots")
	cmd.Flags().BoolVarP(&opts.Create, "create", "c", false, "Create a new test case")

	if err := option.ApplyFlags(&opts, cmd.Flags()); err != nil {
		return nil
	}

	return cmd
}

func runTest(cmd *cobra.Command, opts testOptions) error {
	re, err := recipe.LoadRecipe(opts.RecipePath)
	if err != nil {
		return fmt.Errorf("can not load the recipe: %w", err)
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
			return fmt.Errorf("failed to save recipe: %w", err)
		}

		cmd.Println("Test created")
		return nil
	}

	if len(re.Tests) == 0 {
		cmd.Println("No tests specified")
		return nil
	}

	if opts.UpdateSnapshots {
		for i := range re.Tests {
			test := &re.Tests[i]
			sauce, err := re.Execute(engine.Engine{}, test.Values, recipe.TestID)
			if err != nil {
				return fmt.Errorf("failed to render templates: %w", err)
			}

			test.Files = make(map[string][]byte)
			for filename, file := range sauce.Files {
				test.Files[filename] = file.Content
			}
		}

		err := re.Save(filepath.Dir(opts.RecipePath))
		if err != nil {
			return fmt.Errorf("failed to save recipe: %w", err)
		}

		// TODO: Show which tests were modified
		cmd.Println("Recipe tests updated successfully")
		return nil
	}

	cmd.Printf("Running tests for recipe \"%s\"...\n", re.Name)
	errs := re.RunTests()
	for i, err := range errs {
		var symbol rune
		if err == nil {
			symbol = '✅'
		} else {
			symbol = '❌'
		}

		cmd.Printf("%c: %s\n", symbol, re.Tests[i].Name)
	}

	formattedErrs := make([]error, 0, len(errs))
	for i, err := range errs {
		if err != nil {
			formattedErrs = append(formattedErrs, fmt.Errorf("test %s failed: %v", re.Tests[i].Name, err))
		}
	}

	if len(formattedErrs) > 0 {
		return fmt.Errorf("recipe tests failed: %w", errors.Join(formattedErrs...))
	}

	return nil
}
