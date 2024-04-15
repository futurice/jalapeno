package cli

import (
	"github.com/spf13/cobra"
)

func NewCreateCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use: "create",
	}

	cmd.AddCommand(
		NewCreateRecipeCmd(),
		NewCreateManifestCmd(),
		NewCreateTestCmd(),
	)

	return cmd
}
