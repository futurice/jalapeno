package recipe

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-yaml/yaml"
)

const defaultFileMode os.FileMode = 0700

// Save saves recipe to given destination
func (re *Recipe) Save(dest string) error {
	// TODO: Make sure recipe name is path friendly
	recipeDir := filepath.Join(dest, re.Name)

	// TODO: Override recipe if already exists?
	if _, err := os.Stat(recipeDir); !os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' already exists", dest)
	}

	err := os.Mkdir(recipeDir, defaultFileMode)
	if err != nil {
		return fmt.Errorf("can not create directory %s: %v", recipeDir, err)
	}

	recipeFilepath := filepath.Join(recipeDir, RecipeFileName+YAMLExtension)
	file, err := os.Create(recipeFilepath)
	if err != nil {
		return fmt.Errorf("failed to create recipe file: %w", err)
	}

	encoder := yaml.NewEncoder(file)
	if err := encoder.Encode(re); err != nil {
		return fmt.Errorf("failed to write recipe test to a file: %w", err)
	}

	if err := encoder.Close(); err != nil {
		return fmt.Errorf("failed to close recipe YAML encoder: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("failed to close recipe file: %w", err)
	}

	err = re.saveTemplates(recipeDir)
	if err != nil {
		return fmt.Errorf("can not save recipe templates: %w", err)
	}

	err = re.saveTests(recipeDir)
	if err != nil {
		return fmt.Errorf("can not save recipe tests: %w", err)
	}

	return nil
}

func (re *Recipe) saveTests(dest string) error {
	if len(re.Tests) == 0 {
		return nil
	}

	testDir := filepath.Join(dest, RecipeTestsDirName)

	err := os.Mkdir(testDir, defaultFileMode)
	if err != nil {
		return fmt.Errorf("can not create recipe test directory: %w", err)
	}

	for _, test := range re.Tests {
		file, err := os.Create(filepath.Join(testDir, test.Name+YAMLExtension))
		if err != nil {
			return fmt.Errorf("failed to create rendered recipe file: %w", err)
		}

		encoder := yaml.NewEncoder(file)

		if err := encoder.Encode(test); err != nil {
			return fmt.Errorf("failed to write recipe test to a file: %w", err)
		}

		if err := encoder.Close(); err != nil {
			return fmt.Errorf("failed to close recipe test YAML encoder: %w", err)
		}
		if err := file.Close(); err != nil {
			return fmt.Errorf("failed to close recipe test file: %w", err)
		}
	}

	return nil
}

func (re *Recipe) saveTemplates(dest string) error {
	templateDir := filepath.Join(dest, RecipeTemplatesDirName)
	err := os.Mkdir(templateDir, defaultFileMode)
	if err != nil {
		return fmt.Errorf("can not save templates to the directory: %w", err)
	}

	for relativeFilepath, content := range re.Templates {
		templatePath := filepath.Join(templateDir, relativeFilepath)
		err = os.MkdirAll(filepath.Dir(templatePath), defaultFileMode)
		if err != nil {
			return fmt.Errorf("failed to create rendered recipe file: %w", err)
		}

		err := os.WriteFile(templatePath, content, defaultFileMode)
		if err != nil {
			return fmt.Errorf("failed to create rendered recipe file: %w", err)
		}
	}

	return nil
}

// Save saves sauce to given destination
func (s *Sauce) Save(dest string) error {
	// load all sauces from target dir, because we will either replace
	// a previous rendering of this recipe, or create a new file
	sauces, err := LoadSauces(dest)
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

	if err := os.MkdirAll(filepath.Join(dest, SauceDirName), defaultFileMode); err != nil {
		return fmt.Errorf("failed to create rendered recipe dir: %w", err)
	}
	file, err := os.Create(filepath.Join(dest, SauceDirName, SaucesFileName+YAMLExtension))
	if err != nil {
		return fmt.Errorf("failed to create rendered recipe file: %w", err)
	}
	encoder := yaml.NewEncoder(file)

	for _, sauce := range sauces {
		if err := encoder.Encode(sauce); err != nil {
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
