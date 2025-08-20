package cli

import (
	"github.com/futurice/jalapeno/internal/cli/option"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/spf13/cobra"
)

type listOptions struct {
	option.Common
	option.WorkingDirectory
}

func NewListCmd() *cobra.Command {
	var opts listOptions
	var cmd = &cobra.Command{
		Use:   "list",
		Short: "List sauce(s) in the project",
		Long:  "List installed sauce(s) in the project.",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return option.Parse(&opts)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runList(cmd, opts)
			return errorHandler(cmd, err)
		},
		Example: `# List sauces in the project
jalapeno list`,
	}

	if err := option.ApplyFlags(&opts, cmd.Flags()); err != nil {
		return nil
	}

	return cmd
}

func runList(cmd *cobra.Command, opts listOptions) error {
	sauces, err := recipe.LoadSauces(opts.Dir)
	if err != nil {
		return err
	}

	if len(sauces) == 0 {
		cmd.Println("No sauces found")
		return nil
	}

	cmd.Println("Found sauces:")
	for _, sauce := range sauces {
		cmd.Printf("- name: %s\n  version: %s\n  id: %s\n", sauce.Recipe.Name, sauce.Recipe.Version, sauce.ID)
	}

	return nil
}
