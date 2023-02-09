package main

import (
	"github.com/futurice/jalapeno/cmd/internal/option"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/spf13/cobra"
)

type validateOptions struct {
	TargetPath string
	option.Common
}

func newValidateCmd() *cobra.Command {
	var opts validateOptions
	var cmd = &cobra.Command{
		Use:   "validate RECIPE",
		Short: "Validate a recipe",
		Long:  "", // TODO
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.TargetPath = args[0]
			return option.Parse(&opts)
		},
		Run: func(cmd *cobra.Command, args []string) {
			runValidate(cmd, opts)
		},
		Args: cobra.ExactArgs(1),
	}

	if err := option.ApplyFlags(&opts, cmd.Flags()); err != nil {
		return nil
	}

	return cmd
}

func runValidate(cmd *cobra.Command, opts validateOptions) {
	r, err := recipe.Load(opts.TargetPath)
	if err != nil {
		cmd.PrintErrf("could not load the recipe: %v\n", err)
	}

	err = r.Validate()
	if err != nil {
		cmd.PrintErrf("validation failed: %v\n", err)
	}

	cmd.Println("Validation ok")
}
