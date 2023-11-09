package cli

import (
	"fmt"
	"time"

	"github.com/carlmjohnson/versioninfo"
	"github.com/spf13/cobra"
)

type ExitCodeContextKey struct{}

func NewRootCmd(version string) *cobra.Command {
	// rootCmd represents the base command when called without any subcommands
	var cmd = &cobra.Command{
		Use:          "jalapeno",
		Short:        "Create, manage and share spiced up project templates",
		Long:         "Create, manage and share spiced up project templates.",
		SilenceUsage: true,
	}

	if version != "" {
		cmd.Version = version
	} else {
		cmd.Version = fmt.Sprintf(
			"%s (Built on %s from Git SHA %s)",
			versioninfo.Version,
			versioninfo.Revision,
			versioninfo.LastCommit.Format(time.RFC3339),
		)
	}

	cmd.AddCommand(
		NewCheckCmd(),
		NewCreateCmd(),
		NewEjectCmd(),
		NewExecuteCmd(),
		NewPullCmd(),
		NewPushCmd(),
		NewTestCmd(),
		NewUpgradeCmd(),
		NewValidateCmd(),
		NewWhyCmd(),
	)

	return cmd
}
