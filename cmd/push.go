package main

import (
	"context"
	"strings"

	"github.com/futurice/jalapeno/pkg/recipe"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/spf13/cobra"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"
)

func newPushCmd() *cobra.Command {
	var pushCmd = &cobra.Command{
		Use:   "push",
		Short: "Push a recipe to OCI repository",
		Long:  "",
		Run:   pushFunc,
	}

	return pushCmd
}

func pushFunc(cmd *cobra.Command, args []string) {
	path := args[0]
	targetRef := args[1]
	ctx := context.Background()

	re, err := recipe.Load(path)
	if err != nil {
		cmd.PrintErrf("Error: can't load the recipe: %s\n", err)
		return
	}

	store := file.New("")
	defer store.Close()

	desc, err := store.Add(ctx, re.Name, "application/x.futurice.jalapeno.recipe.v1", path)
	if err != nil {
		cmd.PrintErr(err)
		return
	}

	root, err := oras.Pack(ctx, store, "", []v1.Descriptor{desc}, oras.PackOptions{PackImageManifest: true})
	if err != nil {
		cmd.PrintErr(err)
		return
	}

	err = store.Tag(ctx, root, re.Version)
	if err != nil {
		cmd.PrintErr(err)
		return
	}

	repo, err := remote.NewRepository(targetRef)
	if err != nil {
		cmd.PrintErr(err)
		return
	}

	if strings.Contains(targetRef, "localhost") {
		repo.PlainHTTP = true
	}

	_, err = oras.Copy(ctx, store, re.Version, repo, repo.Reference.Reference, oras.DefaultCopyOptions)
	if err != nil {
		cmd.PrintErr(err)
		return
	}
}
