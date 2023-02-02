package main

import (
	"context"
	"strings"

	"github.com/spf13/cobra"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"
)

func newPullCmd() *cobra.Command {
	var pullCmd = &cobra.Command{
		Use:   "pull",
		Short: "Pull a recipe from OCI repository",
		Long:  "",
		Args:  cobra.ExactArgs(1),
		Run:   pullFunc,
	}

	pullCmd.Flags().StringVarP(&outputBasePath, "output", "o", ".", "Path where the recipe files should be saved")

	return pullCmd
}

func pullFunc(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	targetRef := args[0]

	repo, err := remote.NewRepository(targetRef)
	if err != nil {
		cmd.PrintErr(err)
		return
	}

	if strings.Contains(targetRef, "localhost") {
		repo.PlainHTTP = true
	}

	dst := file.New(outputBasePath)
	_, err = oras.Copy(ctx, repo, repo.Reference.Reference, dst, repo.Reference.Reference, oras.DefaultCopyOptions)
	if err != nil {
		cmd.PrintErr(err)
		return
	}
}
