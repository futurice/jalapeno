package recipeutil

import (
	"errors"
	"os"
	"path/filepath"
)

func SaveFiles(files map[string][]byte, dest string) error {
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		return errors.New("destination path does not exist")
	}

	for name, data := range files {
		path := filepath.Join(dest, name)

		// Create file's parent directories (if not already exist)
		err := os.MkdirAll(filepath.Dir(path), 0700)
		if err != nil {
			return err
		}

		// Create the file
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()

		// Write the data to the file
		_, err = f.Write(data)
		if err != nil {
			return err
		}

		f.Sync()
	}
	return nil
}
