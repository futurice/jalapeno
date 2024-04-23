package cli

import (
	"errors"
	"fmt"

	"github.com/futurice/jalapeno/internal/cli/option"
	"github.com/futurice/jalapeno/pkg/engine"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/recipeutil"
	"github.com/spf13/cobra"
)

type testOptions struct {
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
			err := runTest(cmd, opts)
			return errorHandler(cmd, err)
		},
		Example: `# Run recipe tests
jalapeno test path/to/recipe

# Bootstrap a new test case
jalapeno test path/to/recipe --create

# Update test file snapshots with the current outputs
jalapeno test path/to/recipe --update-snapshots`,
	}

	cmd.Flags().BoolVarP(&opts.UpdateSnapshots, "update-snapshots", "u", false, "Update test file snapshots")

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

	if len(re.Tests) == 0 {
		cmd.Println("No tests specified")
		return nil
	}

	if opts.UpdateSnapshots {
		return updateSnapshots(cmd, opts, re)
	}

	return testRecipe(cmd, re)
}

func testRecipe(cmd *cobra.Command, re *recipe.Recipe) error {
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
		return fmt.Errorf("recipe tests failed: %w", errors.Join(formattedErrs...))
	}

	return nil
}

func updateSnapshots(cmd *cobra.Command, opts testOptions, re *recipe.Recipe) error {
	anyUpdatesFound := false
	for i := range re.Tests {
		test := &re.Tests[i]
		sauce, err := re.Execute(engine.New(), test.Values, recipe.TestID)
		if err != nil {
			return fmt.Errorf("failed to render templates: %w", err)
		}

		fileStatuses := make(map[string]recipeutil.FileStatus, len(sauce.Files))

		for filename, file := range sauce.Files {
			_, found := test.Files[filename]
			if test.IgnoreExtraFiles && !found {
				continue
			}

			if !found {
				fileStatuses[filename] = recipeutil.FileAdded
			} else if file.Checksum == test.Files[filename].Checksum {
				fileStatuses[filename] = recipeutil.FileUnchanged
			} else {
				fileStatuses[filename] = recipeutil.FileModified
			}
		}

		// Update test files
		for filename := range fileStatuses {
			test.Files[filename] = sauce.Files[filename]
		}

		// Remove files which do not exist anymore in the templates
		for filename := range test.Files {
			if _, found := sauce.Files[filename]; !found {
				delete(test.Files, filename)
				fileStatuses[filename] = recipeutil.FileDeleted
			}
		}

		for _, status := range fileStatuses {
			if status != recipeutil.FileUnchanged {
				if !anyUpdatesFound {
					cmd.Print("Updating snapshots for the following tests:\n\n")
					anyUpdatesFound = true
				}
				tree := recipeutil.CreateFileTree(fmt.Sprintf("%s/files", test.Name), fileStatuses)
				cmd.Println(tree)
				break
			}
		}
	}

	if !anyUpdatesFound {
		cmd.Println("No snapshot updates required for any test.")
		return nil
	}

	err := re.Save(opts.RecipePath)
	if err != nil {
		return fmt.Errorf("failed to save recipe: %w", err)
	}

	cmd.Printf("Recipe test snapshots updated %s\n", opts.Colors.Green.Render("successfully!"))
	return nil
}
