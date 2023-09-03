package cli

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/futurice/jalapeno/internal/cli/internal/option"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/spf13/cobra"
)

type whyOptions struct {
	Filepath string
	option.WorkingDirectory
	option.Common
}

func NewWhyCmd() *cobra.Command {
	var opts whyOptions
	var cmd = &cobra.Command{
		Use:   "why FILEPATH",
		Short: "Explains where a file comes from",
		Long:  "Explains where a file comes from in the project, e.g. which file is recipe or user created",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.Filepath = filepath.Clean(args[0])
			return option.Parse(&opts)
		},
		Run: func(cmd *cobra.Command, args []string) {
			runWhy(cmd, opts)
		},
	}

	if err := option.ApplyFlags(&opts, cmd.Flags()); err != nil {
		return nil
	}

	return cmd
}

func runWhy(cmd *cobra.Command, opts whyOptions) {
	// Supporting absolute paths is not trivial, since then we can't use opts.Dir
	// to find the project root directory and we need to travel up the tree to find
	// the project root
	if filepath.IsAbs(opts.Filepath) {
		cmd.PrintErrln("Error: use path relative to the project directory")
		return
	}

	fileinfo, err := os.Stat(filepath.Join(opts.Dir, opts.Filepath))
	if os.IsNotExist(err) {
		cmd.PrintErrf("File \"%s\" does not exist\n", filepath.Join(opts.Dir, opts.Filepath))
		return
	}

	sauces, err := recipe.LoadSauces(opts.Dir)
	if err != nil {
		cmd.PrintErrf("Error: can not load sauces: %s\n", err)
		return
	}

	if len(sauces) == 0 {
		cmd.PrintErrf("Error: \"%s\" is not a project directory", opts.Dir)
		return
	}

	if filepath.Base(opts.Filepath) == ".jalapeno" {
		cmd.Printf("File \"%s\" is created by Jalapeno\n", opts.Filepath)
		return
	}

	for _, sauce := range sauces {
		for file := range sauce.Files {
			if fileinfo.IsDir() {
				if strings.HasPrefix(file, opts.Filepath) {
					cmd.Printf("Directory \"%s\" is created by the recipe \"%s\"\n", opts.Filepath, sauce.Recipe.Name)
					return
				}
			}
			if opts.Filepath == file {
				cmd.Printf("File \"%s\" is created by the recipe \"%s\"\n", opts.Filepath, sauce.Recipe.Name)
				return
			}
		}
	}

	cmd.Printf("File %s is created by the user\n", opts.Filepath)
}
