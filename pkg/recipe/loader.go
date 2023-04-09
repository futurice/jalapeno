package recipe

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-yaml/yaml"
)

const (
	YAMLExtension          = ".yml"
	RecipeFileName         = "recipe"
	RecipeTemplatesDirName = "templates"
	RecipeTestsDirName     = "tests"
	IgnoreFileName         = ".jalapenoignore"
)

// LoadRecipe a recipe from a path. The function does validate the recipe before returning it
func LoadRecipe(path string) (*Recipe, error) {
	rootDir, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	recipeFile := filepath.Join(rootDir, RecipeFileName+YAMLExtension)
	dat, err := os.ReadFile(recipeFile)
	if err != nil {
		return nil, err
	}

	recipe := NewRecipe()
	err = yaml.Unmarshal(dat, recipe)
	if err != nil {
		return nil, err
	}

	recipe.Templates, err = loadTemplates(filepath.Join(rootDir, RecipeTemplatesDirName))
	if err != nil {
		return nil, fmt.Errorf("error when loading recipe templates: %w", err)
	}

	recipe.Tests, err = loadTests(filepath.Join(rootDir, RecipeTestsDirName))
	if err != nil {
		return nil, fmt.Errorf("error when loading recipe tests: %w", err)
	}

	if err := recipe.Validate(); err != nil {
		return nil, err
	}

	return recipe, nil
}

func loadTemplates(recipePath string) (map[string][]byte, error) {
	templates := make(map[string][]byte)

	walk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Continue walking if the path is directory
		if info.IsDir() {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Create a filepath related to the root of the directory
		prefix := fmt.Sprintf("%s%c", recipePath, filepath.Separator)
		name := filepath.ToSlash(strings.TrimPrefix(path, prefix))

		templates[name] = data
		return nil
	}

	err := filepath.Walk(recipePath, walk)
	if err != nil {
		return nil, err
	}

	return templates, nil
}

func loadTests(path string) ([]Test, error) {
	tests := make([]Test, 0)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return tests, nil
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		// Check if not valid test file
		if !strings.HasSuffix(file.Name(), YAMLExtension) || file.IsDir() {
			continue
		}

		test := Test{}
		contents, err := os.ReadFile(filepath.Join(path, file.Name()))
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(contents, &test)
		if err != nil {
			return nil, err
		}

		// If the test does not define the name, get it from filename
		if test.Name == "" {
			test.Name = strings.TrimSuffix(file.Name(), YAMLExtension)
		}

		tests = append(tests, test)
	}

	return tests, nil
}

// Load all sauces from a project directory
func LoadSauces(projectDir string) ([]*Sauce, error) {
	var sauces []*Sauce

	sauceFile := filepath.Join(projectDir, SauceDirName, SaucesFileName+YAMLExtension)
	if _, err := os.Stat(sauceFile); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// treat missing file as empty
			return sauces, nil
		}
		// other errors go boom in os.ReadFile() below
	}
	recipedata, err := os.ReadFile(sauceFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read recipe file: %w", err)
	}

	decoder := yaml.NewDecoder(bytes.NewReader(recipedata))
	for {
		sauce := NewSauce()
		if err := decoder.Decode(&sauce); err != nil {
			if err != io.EOF {
				return nil, fmt.Errorf("failed to decode recipe: %w", err)
			}
			// ran out of recipe file, all yaml documents read
			break
		}
		// read rendered files
		for path, file := range sauce.Files {
			data, err := os.ReadFile(filepath.Join(projectDir, path))
			if err != nil {
				return nil, fmt.Errorf("failed to read rendered file: %w", err)
			}
			file.Content = data
			sauce.Files[path] = file
		}

		if err := sauce.Validate(); err != nil {
			return nil, fmt.Errorf("failed to validate recipe: %w", err)
		}

		sauces = append(sauces, sauce)
	}

	return sauces, nil
}
