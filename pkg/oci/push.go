package oci

import (
	"context"
	"errors"
	"strings"

	"github.com/futurice/jalapeno/pkg/recipe"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
)

var (
	ErrorUnauthorized = errors.New("unauthorized")
)

func PushRecipe(ctx context.Context, path string, opts Repository) error {
	re, err := recipe.LoadRecipe(path)
	if err != nil {
		return err
	}

	repo, err := NewRepository(opts)
	if err != nil {
		return err
	}

	store, err := file.New("")
	if err != nil {
		return err
	}

	defer store.Close()

	desc, err := store.Add(ctx, re.Name, "application/x.futurice.jalapeno.recipe.v1", path)
	if err != nil {
		return err
	}

	root, err := oras.Pack(ctx, store, "", []v1.Descriptor{desc}, oras.PackOptions{PackImageManifest: true})
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
