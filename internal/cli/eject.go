package cli

import (
	"os"
	"path/filepath"

	"github.com/futurice/jalapeno/internal/cli/internal/option"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/spf13/cobra"
)

type ejectOptions struct {
	option.WorkingDirectory
}

func NewEjectCmd() *cobra.Command {
	var opts ejectOptions
	var cmd = &cobra.Command{
		Use:   "eject",
		Short: "Remove all Jalapeno-specific files from a project",
		Long:  "Remove all the files and directories that are for Jalapeno internal use, and leave only the rendered project files.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return option.Parse(&opts)
		},
		Run: func(cmd *cobra.Command, args []string) {
			runEject(cmd, opts)
		},
	}

	if err := option.ApplyFlags(&opts, cmd.Flags()); err != nil {
		return nil
	}

	return cmd
}

func runEject(cmd *cobra.Command, opts ejectOptions) {
	if _, err := os.Stat(opts.Dir); os.IsNotExist(err) {
		cmd.PrintErrln("Error: project path does not exist")
		return
	}

	jalapenoPath := filepath.Join(opts.Dir, recipe.SauceDirName)

	if stat, err := os.Stat(jalapenoPath); os.IsNotExist(err) || !stat.IsDir() {
		cmd.PrintErrf("Error: '%s' is not a Jalapeno project\n", opts.Dir)
		return
	}

	cmd.Printf("Deleting %s...", jalapenoPath)
	err := os.RemoveAll(jalapenoPath)
	if err != nil {
		cmd.PrintErrf("Error: %s", err)
		return
	}

	cmd.Println("\nEjected successfully!")
}
