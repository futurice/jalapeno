package main

import (
	"io"

	"github.com/spf13/cobra"
)

func newRootCmd(out io.Writer, args []string) (*cobra.Command, error) {
	// rootCmd represents the base command when called without any subcommands
	var cmd = &cobra.Command{
		Use:   "jalapeno",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
	examples and usage of using your application. For example:

	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
	}

	cmd.AddCommand(
		newCreateCmd(),
		newUpgradeCmd(),
		newExecuteCmd(),
		newValidateCmd(),
	)

	return cmd, nil
}
