package recipe

import (
	"os"
	"path/filepath"

	"github.com/go-yaml/yaml"
)

// Saves recipe file to given destination
func (re *Recipe) Save(dest string) error {
	out, err := yaml.Marshal(re)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(dest, RecipeFileName), out, 0700)
	if err != nil {
		return err
	}

	return nil
}
