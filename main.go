package main

import (
	"context"

	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"
)

// Before running the code, setup local container registry with
// `docker run -d -p 5001:5000 --restart=always --name registry registry:2`

func main() {
	recipeVersion := "v0.0.1" // NOTE: This should come from recipe file

	ctx := context.Background()

	store := file.New("")
	defer store.Close()

	desc1, err := store.Add(ctx, "recipe.yml", "", "./examples/minimal/recipe.yml")
	check(err)
	desc2, err := store.Add(ctx, "templates/README.md", "", "./examples/minimal/templates/README.md")
	check(err)

	root, err := oras.Pack(ctx, store, "", []v1.Descriptor{desc1, desc2}, oras.PackOptions{PackImageManifest: true})
	check(err)

	err = store.Tag(ctx, root, recipeVersion)
	check(err)

	reg, err := remote.NewRepository("localhost:5001/minimal")
	check(err)

	reg.PlainHTTP = true

	// Push the artifact to local registry
	_, err = oras.Copy(ctx, store, recipeVersion, reg, "", oras.DefaultCopyOptions)
	check(err)

	// Download the artifact from local registry
	dst := file.New("./output")
	_, err = oras.Copy(ctx, reg, recipeVersion, dst, "", oras.DefaultCopyOptions)
	check(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
