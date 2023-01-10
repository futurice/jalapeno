package main

import (
	"context"

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
		Short: "",
		Long:  "",
		Run:   pushFunc,
	}

	return pushCmd
}

func pushFunc(cmd *cobra.Command, args []string) {
	path := args[0]
	ctx := context.Background()

	re, err := recipe.Load(path)
	if err != nil {
		cmd.PrintErrf("Error: can't load the recipe: %s\n", err)
		return
	}

	store := file.New("")
	defer store.Close()

	desc, err := store.Add(ctx, "recipe", "", path)
	check(err)

	root, err := oras.Pack(ctx, store, "", []v1.Descriptor{desc}, oras.PackOptions{PackImageManifest: true})
	check(err)

	err = store.Tag(ctx, root, re.Version)
	check(err)

	repo, err := remote.NewRepository(args[1])
	check(err)

	repo.PlainHTTP = true

	_, err = oras.Copy(ctx, store, re.Version, repo, repo.Reference.Reference, oras.DefaultCopyOptions)
	check(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
