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
	Create          bool
	RecipePath      string
	UpdateSnapshots bool

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
		anyUpdatesFound := false
		for i := range re.Tests {
			test := &re.Tests[i]
			fileUpdatesFound := make(map[string]recipe.File)
			sauce, err := re.Execute(engine.New(), test.Values, recipe.TestID)
			if err != nil {
				return fmt.Errorf("failed to render templates: %w", err)
			}

			for filename, file := range sauce.Files {
				if _, found := test.Files[filename]; test.IgnoreExtraFiles && !found {
					continue
				}

				if file.Checksum != test.Files[filename].Checksum {
					fileUpdatesFound[filename] = file
				}
			}

			// Update test files
			for filename, update := range fileUpdatesFound {
				test.Files[filename] = update
			}

			// Remove files which do not exist anymore in the templates
			for filename := range test.Files {
				if _, found := sauce.Files[filename]; !found {
					delete(test.Files, filename)
				}
			}

			if len(fileUpdatesFound) > 0 {
				if !anyUpdatesFound {
					cmd.Print("The following files have been updated:\n\n")
					anyUpdatesFound = true
				}

				cmd.Print(
					recipeutil.CreateFileTree(
						fmt.Sprintf("%s/files", test.Name),
						fileUpdatesFound,
					),
				)
			}
		}

		if !anyUpdatesFound {
			cmd.Println("No snapshot updates for any test")
			return nil
		}

		err := re.Save(filepath.Dir(opts.RecipePath))
		if err != nil {
			return fmt.Errorf("failed to save recipe: %w", err)
		}

		cmd.Printf("\nRecipe test snapshots updated %s\n", opts.Colors.Green.Render("successfully!"))
		return nil
	}

	cmd.Printf("Running tests for recipe \"%s\"...\n\n", re.Name)
	errs := re.RunTests()
	for i, err := range errs {
		symbol := '✅'
		if err != nil {
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
		cmd.Println()
		return fmt.Errorf("recipe tests failed: %w", errors.Join(formattedErrs...))
	}

	return nil
}
