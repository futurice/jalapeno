package main

import (
	"io"

	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	cmd, err := newRootCmd(os.Stdout, os.Args[1:])
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}

	if err = cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd(out io.Writer, args []string) (*cobra.Command, error) {
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
	)

	return cmd, nil
}
