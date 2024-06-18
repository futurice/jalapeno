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

const LatestTag = "latest"

var (
	ErrorUnauthorized = errors.New("unauthorized")
)

func PushRecipe(ctx context.Context, path string, opts oci.Repository, replaceLatest bool) error {
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

	// TODO: Add each file separately so media types can be set correctly. Also, the recipe directory
	//       could contain additional files which are not related directly to the recipe.
	mediaType := "inode/directory"
	desc, err := store.Add(ctx, re.Name, mediaType, path)
	if err != nil {
		return err
	}
	fileDescriptors := []v1.Descriptor{desc}

	root, err := oras.PackManifest(
		ctx,
		store,
		oras.PackManifestVersion1_1,
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

	destTags := []string{re.Version}
	if replaceLatest {
		destTags = append(destTags, LatestTag)
	}

	for _, tag := range destTags {
		_, err = oras.Copy(ctx, store, re.Version, repo, tag, oras.DefaultCopyOptions)
		if err != nil {
			if strings.Contains(err.Error(), "credential required") {
				return ErrorUnauthorized
			} else {
				return err
			}
		}
	}

	return nil
}
