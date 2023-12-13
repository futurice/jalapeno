package recipe

import (
	"context"
	"errors"
	"strings"

	"github.com/futurice/jalapeno/pkg/oci"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
)

var (
	ErrorUnauthorized = errors.New("unauthorized")
)

func PushRecipe(ctx context.Context, path string, opts oci.Repository) error {
	re, err := LoadRecipe(path)
	if err != nil {
		return err
	}

	repo, err := oci.NewRepository(opts)
	if err != nil {
		return err
	}

	store, err := file.New("")
	if err != nil {
		return err
	}

	defer store.Close()

	// TODO: Add each file separately so media types can be set correctly
	mediaType := "inode/directory"
	desc, err := store.Add(ctx, re.Name, mediaType, path)
	if err != nil {
		return err
	}
	fileDescriptors := []v1.Descriptor{desc}

	root, err := oras.PackManifest(
		ctx,
		store,
		oras.PackManifestVersion1_1_RC4,
		"application/x.futurice.jalapeno.recipe.v1",
		oras.PackManifestOptions{
			Layers: fileDescriptors,
		})
	if err != nil {
		return err
	}

	err = store.Tag(ctx, root, re.Version)
	if err != nil {
		return err
	}

	_, err = oras.Copy(ctx, store, re.Version, repo, re.Version, oras.DefaultCopyOptions)
	if err != nil {
		if strings.Contains(err.Error(), "credential required") {
			return ErrorUnauthorized
		} else {
			return err
		}
	}

	return nil
}
