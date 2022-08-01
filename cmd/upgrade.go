package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newUpgradeCmd() *cobra.Command {
	// upgradeCmd represents the upgrade command
	var upgradeCmd = &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade recipe in a project",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("upgrade called")
		},
	}

	return upgradeCmd
}
