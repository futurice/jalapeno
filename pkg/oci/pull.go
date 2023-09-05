package oci

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/futurice/jalapeno/pkg/recipe"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
)

func PullRecipe(ctx context.Context, opts Repository) (*recipe.Recipe, error) {
	dir, err := os.MkdirTemp("", "jalapeno-remote-recipe")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)

	err = SaveRemoteRecipe(ctx, dir, opts)
	if err != nil {
		return nil, err
	}

	// We don't know the recipe's real name at this point, so we need to check the directory
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	if len(entries) != 1 || !entries[0].IsDir() {
		return nil, errors.New("after pulling remote recipe the directory did not contain assumed directory")
	}

	recipeName := entries[0].Name()

	re, err := recipe.LoadRecipe(filepath.Join(dir, recipeName))
	if err != nil {
		return nil, err
	}

	return re, nil
}

// SaveRemoteRecipe pulls a recipe from repository and saves it to dest directory
func SaveRemoteRecipe(ctx context.Context, dest string, opts Repository) error {
	repo, err := NewRepository(opts)
	if err != nil {
		return err
	}

	dst, err := file.New(dest)
	if err != nil {
		return err
	}

	_, err = oras.Copy(
		ctx,
		repo,
		repo.Reference.Reference,
		dst,
		repo.Reference.Reference,
		oras.DefaultCopyOptions,
	)
	if err != nil {
		return err
	}

	return nil
}
