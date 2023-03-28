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
	RecipeFileName         = "recipe.yml"
	RecipeTemplatesDirName = "templates"
	RecipeTestsDirName     = "tests"
	RenderedRecipeDirName  = ".jalapeno"
	IgnoreFileName         = ".jalapenoignore"
)

// Load a recipe from a path. The function does validate the recipe before returning it
func Load(path string) (*Recipe, error) {
	rootDir, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	recipeFile := filepath.Join(rootDir, RecipeFileName)
	dat, err := os.ReadFile(recipeFile)
	if err != nil {
		return nil, err
	}

	recipe := new()
	err = yaml.Unmarshal(dat, recipe)
	if err != nil {
		return nil, err
	}

	recipe.Templates, err = loadTemplates(filepath.Join(rootDir, RecipeTemplatesDirName))
	if err != nil {
		return nil, err
	}

	recipe.tests, err = loadTests(filepath.Join(rootDir, RecipeTestsDirName))
	if err != nil {
		return nil, err
	}

	if err := recipe.Validate(); err != nil {
		return nil, err
	}

	return recipe, nil
}

// Load recipes which already have been rendered. Always loads
// all recipes.
func LoadRendered(projectDir string) ([]*Recipe, error) {
	var recipes []*Recipe

	recipeFile := filepath.Join(projectDir, RenderedRecipeDirName, RecipeFileName)
	if _, err := os.Stat(recipeFile); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// treat missing file as empty
			return recipes, nil
		}
		// other errors go boom in os.ReadFile() below
	}
	recipedata, err := os.ReadFile(recipeFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read recipe file: %w", err)
	}

	decoder := yaml.NewDecoder(bytes.NewReader(recipedata))
	for {
		recipe := new()
		if err := decoder.Decode(&recipe); err != nil {
			if err != io.EOF {
				return nil, fmt.Errorf("failed to decode recipe: %w", err)
			}
			// ran out of recipe file, all yaml documents read
			break
		}
		// read rendered files
		for path, file := range recipe.Files {
			data, err := os.ReadFile(filepath.Join(projectDir, path))
			if err != nil {
				return nil, fmt.Errorf("failed to read rendered file: %w", err)
			}
			file.Content = data
			recipe.Files[path] = file
		}

		if err := recipe.Validate(); err != nil {
			return nil, fmt.Errorf("failed to validate recipe: %w", err)
		}

		recipes = append(recipes, recipe)
	}

	return recipes, nil
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

	// TODO
	tests = append(tests, Test{
		Name:   "case-1",
		Values: VariableValues{},
		Files: map[string][]byte{
			"README.md": []byte("# minimal recipes, version v0.0.1"),
		},
	})

	return tests, nil
}
