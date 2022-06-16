package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newUpdateCmd() *cobra.Command {
	// updateCmd represents the update command
	var updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update the CLI",
		Long:  "Update the CLI to latest version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("update called")
		},
	}

	return updateCmd
}
