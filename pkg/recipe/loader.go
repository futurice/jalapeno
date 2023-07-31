package recipe

import (
	"errors"
	"fmt"
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

var (
	ErrSauceNotFound = errors.New("sauce not found")
)

func LoadRecipe(url string) (*Recipe, error) {
	// TODO: Check the protocol and choose the correct loader
	return loadRecipeFromDir(url)
}

// loadRecipeFromDir loads a recipe from a file path
func loadRecipeFromDir(path string) (*Recipe, error) {
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
