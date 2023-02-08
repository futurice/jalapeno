package main

import (
	"context"
	"strings"

	"github.com/futurice/jalapeno/internal/option"
	"github.com/spf13/cobra"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"
)

type pullOptions struct {
	TargetRef string

	option.Output
	option.Remote
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

	repo, err := remote.NewRepository(opts.TargetRef)
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	if strings.Contains(opts.TargetRef, "localhost") {
		repo.PlainHTTP = true
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
