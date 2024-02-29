package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/futurice/jalapeno/internal/cli/option"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/spf13/cobra"
)

type ejectOptions struct {
	option.Common
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
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runEject(cmd, opts)
			return errorHandler(cmd, err)
		},
		Example: `jalapeno eject`,
	}

	if err := option.ApplyFlags(&opts, cmd.Flags()); err != nil {
		return nil
	}

	return cmd
}

func runEject(cmd *cobra.Command, opts ejectOptions) error {
	if _, err := os.Stat(opts.Dir); os.IsNotExist(err) {
		return errors.New("project path does not exist")
	}

	jalapenoPath := filepath.Join(opts.Dir, recipe.SauceDirName)

	if stat, err := os.Stat(jalapenoPath); os.IsNotExist(err) || !stat.IsDir() {
		return fmt.Errorf("'%s' is not a Jalapeno project", opts.Dir)
	}

	cmd.Printf("Deleting %s...\n", jalapenoPath)
	err := os.RemoveAll(jalapenoPath)
	if err != nil {
		return err
	}

	cmd.Printf("Ejected %s\n", opts.Colors.Green.Render("successfully!"))
	return nil
}
