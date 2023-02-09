package main

import (
	"os"
	"path/filepath"

	"github.com/futurice/jalapeno/internal/option"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/spf13/cobra"
)

type ejectOptions struct {
	ProjectPath string
}

func newEjectCmd() *cobra.Command {
	var opts ejectOptions
	var cmd = &cobra.Command{
		Use:   "eject (PATH)",
		Short: "Remove all Jalapeno-specific files",
		Long:  "Remove all the files and directories that are for Jalapeno internal use, and leave only the rendered project files.",
		Args:  cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) >= 0 {
				opts.ProjectPath = args[0]
			} else {
				opts.ProjectPath = "."
			}
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
	if _, err := os.Stat(opts.ProjectPath); os.IsNotExist(err) {
		cmd.PrintErrln("project path does not exist")
		return
	}

	jalapenoPath := filepath.Join(opts.ProjectPath, recipe.RenderedRecipeDirName)

	if stat, err := os.Stat(jalapenoPath); os.IsNotExist(err) || !stat.IsDir() {
		cmd.PrintErrf("'%s' is not a Jalapeno project\n", opts.ProjectPath)
		return
	}

	cmd.Printf("Deleting %s...", jalapenoPath)
	err := os.RemoveAll(jalapenoPath)
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	cmd.Println("\nEjected successfully!")
}
