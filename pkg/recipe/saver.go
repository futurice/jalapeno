package recipe

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-yaml/yaml"
)

func (s *Recipe) Save(dest string) error {
	// TODO
	return nil
}

// Saves sauce to given destination
func (s *Sauce) Save(dest string) error {
	// load all sauces from target dir, because we will either replace
	// a previous rendering of this recipe, or create a new file
	sauces, err := LoadSauce(dest)
	if err != nil {
		return err
	}
	added := false
	for i, prev := range sauces {
		if s.Recipe.Name == prev.Recipe.Name {
			// found by name
			sauces[i] = s
			added = true
			break
		}
	}
	if !added {
		// we hit the end, append
		sauces = append(sauces, s)
	}

	if err := os.MkdirAll(filepath.Join(dest, SauceDirName), 0755); err != nil {
		return fmt.Errorf("failed to create rendered recipe dir: %w", err)
	}
	file, err := os.Create(filepath.Join(dest, SauceDirName, SauceFileName+YAMLExtension))
	if err != nil {
		return fmt.Errorf("failed to create rendered recipe file: %w", err)
	}
	encoder := yaml.NewEncoder(file)

	for _, recipe := range sauces {
		if err := encoder.Encode(recipe); err != nil {
			return fmt.Errorf("failed to write recipes: %w", err)
		}
	}
	if err := encoder.Close(); err != nil {
		return fmt.Errorf("failed to close recipe file: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("failed to close recipe file: %w", err)
	}

	return nil
}
