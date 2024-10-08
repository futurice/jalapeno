package cli

import (
	"fmt"

	"github.com/futurice/jalapeno/internal/cli/option"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/recipeutil"
	"github.com/futurice/jalapeno/pkg/ui/colors"
	"github.com/spf13/cobra"
)

type createTestOptions struct {
	option.Common
	option.WorkingDirectory
}

func NewCreateTestCmd() *cobra.Command {
	var opts createTestOptions
	var cmd = &cobra.Command{
		Use:     "test",
		Short:   "Create a recipe test",
		Example: `jalapeno create test`,
		Args:    cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return option.Parse(&opts)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runCreateTest(cmd, opts)
			return errorHandler(cmd, err)
		},
	}

	if err := option.ApplyFlags(&opts, cmd.Flags()); err != nil {
		return nil
	}

	return cmd
}

func runCreateTest(cmd *cobra.Command, opts createTestOptions) error {
	re, err := recipe.LoadRecipe(opts.Dir)
	if err != nil {
		return fmt.Errorf("can not load the recipe: %w", err)
	}

	testCaseName := generateTestCaseName(re.Tests)
	test := recipeutil.CreateExampleTest(testCaseName)
	re.Tests = append(re.Tests, test)

	err = re.Save(opts.Dir)
	if err != nil {
		return fmt.Errorf("failed to save recipe: %w", err)
	}

	cmd.Printf(
		"Test '%s' created %s\n\n",
		test.Name,
		colors.Green.Render("successfully!"),
	)

	fmt.Printf("Following files were created: \n%s", recipeutil.CreateFileTree(opts.Dir, map[string]recipeutil.FileStatus{
		fmt.Sprintf("tests/%s/test.yml", test.Name): recipeutil.FileAdded,
		fmt.Sprintf("tests/%s/files/", test.Name):   recipeutil.FileAdded,
	}))

	return nil
}

func generateTestCaseName(tests []recipe.Test) string {
	defaultName := "example"

	testNameExists := func(name string, tests []recipe.Test) bool {
		for _, t := range tests {
			if t.Name == name {
				return true
			}
		}
		return false
	}

	if !testNameExists(defaultName, tests) {
		return defaultName
	}

	for i := 1; ; i++ {
		name := fmt.Sprintf("%s_%d", defaultName, i)
		if !testNameExists(name, tests) {
			return name
		}
	}
}
