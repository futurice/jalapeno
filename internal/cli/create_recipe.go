package cli

import (
	"fmt"
	"path/filepath"

	"github.com/futurice/jalapeno/internal/cli/option"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/recipeutil"
	"github.com/futurice/jalapeno/pkg/ui/colors"
	"github.com/spf13/cobra"
)

type createRecipeOptions struct {
	RecipeName string

	option.Common
	option.WorkingDirectory
}

func NewCreateRecipeCmd() *cobra.Command {
	var opts createRecipeOptions
	var cmd = &cobra.Command{
		Use:   "recipe RECIPE_NAME",
		Short: "Create a new recipe",
		Long: fmt.Sprintf(`Create a new recipe with the given name.

The following files will be created:
%[1]s
my-recipe
├── %[2]s
├── %[3]s
│   └── README.md
└── %[4]s
    └── defaults
        ├── %[5]s
        └── %[6]s
            └── README.md
%[1]s`,
			"```",
			recipe.RecipeFileName+recipe.YAMLExtension,
			recipe.RecipeTemplatesDirName,
			recipe.RecipeTestsDirName,
			recipe.RecipeTestMetaFileName+recipe.YAMLExtension,
			recipe.RecipeTestFilesDirName,
		),
		Example: `jalapeno create recipe my-recipe`,
		Args:    cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.RecipeName = args[0]
			return option.Parse(&opts)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runCreateRecipe(cmd, opts)
			return errorHandler(cmd, err)
		},
	}

	if err := option.ApplyFlags(&opts, cmd.Flags()); err != nil {
		return nil
	}

	return cmd
}

func runCreateRecipe(cmd *cobra.Command, opts createRecipeOptions) error {
	re := recipeutil.CreateExampleRecipe(opts.RecipeName)

	err := re.Save(filepath.Join(opts.Dir, opts.RecipeName))
	if err != nil {
		return fmt.Errorf("can not save recipe to the directory: %w", err)
	}

	cmd.Printf(
		"Recipe '%s' created %s\n",
		opts.RecipeName,
		colors.Green.Render("successfully!"),
	)

	return nil
}
