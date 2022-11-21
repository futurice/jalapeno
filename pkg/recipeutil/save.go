package recipeutil

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/futurice/jalapeno/pkg/recipe"
)

func SaveFiles(files map[string]recipe.File, dest string) error {
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		return errors.New("destination path does not exist")
	}

	for path, file := range files {
		destPath := filepath.Join(dest, path)

		// Create file's parent directories (if not already exist)
		err := os.MkdirAll(filepath.Dir(destPath), 0700)
		if err != nil {
			return err
		}

		// Create the file
		f, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer f.Close()

		// Write the data to the file
		_, err = f.Write(file.Content)
		if err != nil {
			return err
		}

		err = f.Sync()
		if err != nil {
			return err
		}
	}
	return nil
}
