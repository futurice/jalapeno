package cli

import (
	"fmt"

	"github.com/futurice/jalapeno/internal/cli/internal/option"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/spf13/cobra"
)

type createOptions struct {
	RecipeName string
	option.Output
	option.Common
}

func NewCreateCmd() *cobra.Command {
	var opts createOptions
	var cmd = &cobra.Command{
		Use:   "create RECIPE_NAME",
		Short: "Create a new recipe",
		Long: fmt.Sprintf(`TODO

%[1]s
foo/
  ├── recipe.yml
  ├── templates/
  ├──── README.md
%[1]s`, "```"),
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
	re := createExampleRecipe(opts.RecipeName)

	err := re.Validate()
	if err != nil {
		cmd.PrintErrln("Internal error: placeholder recipe is not valid")
		return
	}

	err = re.Save(opts.OutputPath)
	if err != nil {
		cmd.PrintErrf("Error: can not save recipe to the directory: %v\n", err)
		return
	}
}

func createExampleRecipe(name string) *recipe.Recipe {
	r := recipe.NewRecipe()
	r.Metadata.Name = name
	r.Metadata.Version = "v0.0.0"
	r.Metadata.Description = "Description about what the recipe is used for and what it contains. For example tech stack, cloud environments, tools"
	r.Variables = []recipe.Variable{
		{Name: "MY_VAR", Default: "Hello World!"},
	}
	r.Templates = map[string][]byte{
		"README.md": []byte("{{ .Variables.MY_VAR }}"),
	}

	return r
}
