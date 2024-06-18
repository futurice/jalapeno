package cli

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/carlmjohnson/versioninfo"
	uiutil "github.com/futurice/jalapeno/pkg/ui/util"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	// https://goreleaser.com/cookbooks/using-main.version/
	version string

	ColorRed   = lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4136"))
	ColorGreen = lipgloss.NewStyle().Foreground(lipgloss.Color("#26A568"))
)

type ExitCodeContextKey struct{}

// Execute runs the command and returns the exit code
func Execute(cmd *cobra.Command) int {
	err := cmd.ExecuteContext(context.Background())
	exitCode, isExitCodeSet := cmd.Context().Value(ExitCodeContextKey{}).(int)
	if isExitCodeSet {
		return exitCode
	}

	if err == nil {
		return 0
	} else {
		return 1
	}
}

func NewRootCmd() *cobra.Command {
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

func errorHandler(cmd *cobra.Command, err error) error {
	if err == nil {
		return nil
	}

	// Print empty line before error message
	cmd.Println()

	// If the error is a user abort, don't print the error message
	if errors.Is(err, uiutil.ErrUserAborted) {
		cmd.Println("User aborted")
		return nil
	}

	// Color the error message
	cmd.SetErrPrefix(ColorRed.Render("Error:"))
	return errors.New(ColorRed.Render(err.Error()))
}
