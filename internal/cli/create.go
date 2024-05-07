package cli

import (
	"github.com/spf13/cobra"
)

func NewCreateCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new recipe, manifest or test",
		Long:  "Create a new recipe, manifest or test.",
	}

	cmd.AddCommand(
		NewCreateRecipeCmd(),
		NewCreateManifestCmd(),
		NewCreateTestCmd(),
	)

	return cmd
}
