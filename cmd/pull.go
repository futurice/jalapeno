package main

import (
	"context"
	"strings"

	"github.com/futurice/jalapeno/cmd/internal/option"
	"github.com/spf13/cobra"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
)

type pullOptions struct {
	TargetRef string

	option.Output
	option.Repository
	option.Common
}

func newPullCmd() *cobra.Command {
	var opts pullOptions
	var cmd = &cobra.Command{
		Use:   "pull URL",
		Short: "Pull a recipe from OCI repository",
		Long:  "",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.TargetRef = args[0]
			return option.Parse(&opts)
		},
		Run: func(cmd *cobra.Command, args []string) {
			runPull(cmd, opts)
		},
	}

	if err := option.ApplyFlags(&opts, cmd.Flags()); err != nil {
		return nil
	}

	return cmd
}

func runPull(cmd *cobra.Command, opts pullOptions) {
	ctx := context.Background()

	repo, err := opts.NewRepository(opts.TargetRef, opts.Common)
	if err != nil {
		cmd.PrintErrf("Error: %s", err)
		return
	}

	dst, err := file.New(opts.OutputPath)
	if err != nil {
		cmd.PrintErrf("Error: %s", err)
		return
	}

	_, err = oras.Copy(ctx, repo, repo.Reference.Reference, dst, repo.Reference.Reference, oras.DefaultCopyOptions)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			cmd.PrintErrln("Error: recipe not found") // TODO: Give more descriptive error message
		} else {
			cmd.PrintErrf("Error: unexpected error happened: %s", err)
		}
		return
	}

	cmd.Println("Recipe pulled successfully!")
}
