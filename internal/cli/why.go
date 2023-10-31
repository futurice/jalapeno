package cli

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/futurice/jalapeno/internal/cli/option"
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
		Long:  "Explains where a file comes from in the project, e.g. is the file create by a recipe or user",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.Filepath = filepath.Clean(args[0])
			return option.Parse(&opts)
		},
		Run: func(cmd *cobra.Command, args []string) {
			runWhy(cmd, opts)
		},
		Example: `jalapeno why path/to/file`,
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
		cmd.PrintErrf("File '%s' does not exist\n", filepath.Join(opts.Dir, opts.Filepath))
		return
	}

	sauces, err := recipe.LoadSauces(opts.Dir)
	if err != nil {
		cmd.PrintErrf("Error: can not load sauces: %s\n", err)
		return
	}

	if len(sauces) == 0 {
		cmd.PrintErrf("Error: '%s' is not a project directory\n", opts.Dir)
		return
	}

	if opts.Filepath == recipe.SauceDirName {
		cmd.Printf("Directory '%s' is created by Jalapeno\n", opts.Filepath)
		return

	}

	if strings.Split(opts.Filepath, string(filepath.Separator))[0] == recipe.SauceDirName {
		cmd.Printf("File '%s' is created by Jalapeno\n", opts.Filepath)
		return
	}

	for _, sauce := range sauces {
		for file := range sauce.Files {
			if fileinfo.IsDir() {
				if strings.HasPrefix(file, opts.Filepath) {
					cmd.Printf("Directory '%s' is created by the recipe '%s'\n", opts.Filepath, sauce.Recipe.Name)
					return
				}
			}
			if opts.Filepath == file {
				// TODO: Check if the file is modified by the user by comparing hashes
				cmd.Printf("File '%s' is created by the recipe '%s'\n", opts.Filepath, sauce.Recipe.Name)
				return
			}
		}
	}

	cmd.Printf("File %s is created by the user\n", opts.Filepath)
}
