package main

import (
	"context"
	"strings"

	"github.com/futurice/jalapeno/internal/option"
	"github.com/futurice/jalapeno/pkg/recipe"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/spf13/cobra"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"
)

type pushOptions struct {
	RecipePath string
	TargetRef  string
	option.Remote
}

func newPushCmd() *cobra.Command {
	var opts pushOptions
	var cmd = &cobra.Command{
		Use:   "push RECIPE_PATH TARGET",
		Short: "Push a recipe to OCI repository",
		Long:  "",
		Args:  cobra.ExactArgs(2),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.RecipePath = args[0]
			opts.TargetRef = args[1]
			return option.Parse(&opts)
		},
		Run: func(cmd *cobra.Command, args []string) {
			runPush(cmd, opts)
		},
	}

	if err := option.ApplyFlags(&opts, cmd.Flags()); err != nil {
		return nil
	}

	return cmd
}

func runPush(cmd *cobra.Command, opts pushOptions) {
	ctx := context.Background()

	re, err := recipe.Load(opts.RecipePath)
	if err != nil {
		cmd.PrintErrf("Error: can't load the recipe: %s\n", err)
		return
	}

	store, err := file.New("")
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	defer store.Close()

	desc, err := store.Add(ctx, re.Name, "application/x.futurice.jalapeno.recipe.v1", opts.RecipePath)
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	root, err := oras.Pack(ctx, store, "", []v1.Descriptor{desc}, oras.PackOptions{PackImageManifest: true})
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	err = store.Tag(ctx, root, re.Version)
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	repo, err := remote.NewRepository(opts.TargetRef)
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	// https://github.com/oras-project/oras/blob/a98931bc68a0e7a5de6c51df8e5aa09aadad3057/cmd/oras/internal/option/remote.go#L181
	if strings.Contains(opts.TargetRef, "localhost") {
		repo.PlainHTTP = true
	}

	_, err = oras.Copy(ctx, store, re.Version, repo, repo.Reference.Reference, oras.DefaultCopyOptions)
	if err != nil {
		cmd.PrintErrln(err)
		return
	}
}
