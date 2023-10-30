package cli

import (
	"fmt"

	"github.com/futurice/jalapeno/internal/cli/option"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/recipeutil"
	"github.com/spf13/cobra"
)

type createOptions struct {
	RecipeName string

	option.Common
	option.WorkingDirectory
}

func NewCreateCmd() *cobra.Command {
	var opts createOptions
	var cmd = &cobra.Command{
		Use:   "create RECIPE_NAME",
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
		Example: `jalapeno create my-recipe`,
		Args:    cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.RecipeName = args[0]
			return option.Parse(&opts)
		},
		Run: func(cmd *cobra.Command, args []string) {
			runCreate(cmd, opts)
		},
	}

	if err := option.ApplyFlags(&opts, cmd.Flags()); err != nil {
		return nil
	}

	return cmd
}

func runCreate(cmd *cobra.Command, opts createOptions) {
	re := recipeutil.CreateExampleRecipe(opts.RecipeName)

	err := re.Validate()
	if err != nil {
		cmd.PrintErrln("Internal error: placeholder recipe is not valid")
		return
	}

	err = re.Save(opts.Dir)
	if err != nil {
		cmd.PrintErrf("Error: can not save recipe to the directory: %v\n", err)
		return
	}

	cmd.Printf("Recipe '%s' created!\n", opts.RecipeName)
}
