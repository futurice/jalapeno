package cli

import (
	"fmt"
	"time"

	"github.com/carlmjohnson/versioninfo"
	"github.com/spf13/cobra"
)

func NewRootCmd() (*cobra.Command, error) {
	// rootCmd represents the base command when called without any subcommands
	var cmd = &cobra.Command{
		Use:   "jalapeno",
		Short: "Create, manage and share spiced up project templates",
		Long:  "",
	}

	cmd.Version = fmt.Sprintf(
		"%s (Built on %s from Git SHA %s)",
		versioninfo.Version,
		versioninfo.Revision,
		versioninfo.LastCommit.Format(time.RFC3339),
	)

	cmd.AddCommand(
		NewCreateCmd(),
		NewUpgradeCmd(),
		NewExecuteCmd(),
		NewValidateCmd(),
		NewEjectCmd(),
		NewPushCmd(),
		NewPullCmd(),
		NewTestCmd(),
		NewCheckCmd(),
	)

	return cmd, nil
}
