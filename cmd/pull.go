package main

import (
	"context"

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
		Use:   "pull",
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
		cmd.PrintErrln(err)
		return
	}

	dst, err := file.New(opts.OutputPath)
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	_, err = oras.Copy(ctx, repo, repo.Reference.Reference, dst, repo.Reference.Reference, oras.DefaultCopyOptions)
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	cmd.Println("Recipe pulled successfully!")
}
