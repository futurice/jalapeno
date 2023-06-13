package cli

import (
	"io"

	"github.com/spf13/cobra"
)

func NewRootCmd(out io.Writer, args []string) (*cobra.Command, error) {
	// rootCmd represents the base command when called without any subcommands
	var cmd = &cobra.Command{
		Use:   "jalapeno",
		Short: "Create, manage and share spiced up project templates",
		Long:  "",
	}

	cmd.AddCommand(
		newCreateCmd(),
		newUpgradeCmd(),
		newExecuteCmd(),
		newValidateCmd(),
		newEjectCmd(),
		newPushCmd(),
		newPullCmd(),
		newTestCmd(),
	)

	return cmd, nil
}
