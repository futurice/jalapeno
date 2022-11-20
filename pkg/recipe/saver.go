package recipe

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-yaml/yaml"
)

// Saves recipe file to given destination
func (re *Recipe) Save(dest string) error {
	// load all recipes from target dir, because we will either replace
	// a previous rendering of this recipe, or create a new file
	recipes, err := LoadRendered(dest)
	if err != nil {
		return err
	}
	added := false
	for i, prev := range recipes {
		if re.Name == prev.Name {
			// found by name
			recipes[i] = *re
			added = true
			break
		}
	}
	if !added {
		// we hit the end, append
		recipes = append(recipes, *re)
	}

	if err := os.MkdirAll(filepath.Join(dest, RenderedRecipeDirName), 0755); err != nil {
		return fmt.Errorf("failed to create rendered recipe dir: %w", err)
	}
	file, err := os.Create(filepath.Join(dest, RenderedRecipeDirName, RecipeFileName))
	if err != nil {
		return fmt.Errorf("failed to create rendered recipe file: %w", err)
	}
	encoder := yaml.NewEncoder(file)

	for _, recipe := range recipes {
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
