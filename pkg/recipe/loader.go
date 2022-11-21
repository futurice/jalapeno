package recipe

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-yaml/yaml"
)

const (
	RecipeFileName         = "recipe.yml"
	RecipeTemplatesDirName = "templates"
	RenderedRecipeDirName  = ".jalapeno"
)

// Load a recipe from its source.
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

	recipe := &Recipe{}
	err = yaml.Unmarshal(dat, recipe)
	if err != nil {
		return nil, err
	}

	if err := recipe.Validate(); err != nil {
		return nil, err
	}

	templates := make(map[string][]byte)
	templatesDir := filepath.Join(rootDir, RecipeTemplatesDirName)

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

		prefix := filepath.Join(rootDir, RecipeTemplatesDirName)
		prefix += string(filepath.Separator)
		name := filepath.ToSlash(strings.TrimPrefix(path, prefix))

		templates[name] = data
		return nil
	}

	err = filepath.Walk(templatesDir, walk)
	if err != nil {
		return recipe, err
	}

	recipe.Templates = templates

	return recipe, nil
}

// Load recipes which already have been rendered. Always loads
// all recipes.
func LoadRendered(projectDir string) ([]Recipe, error) {
	var recipes []Recipe

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
		recipe := Recipe{}
		if err := decoder.Decode(&recipe); err != nil {
			if err != io.EOF {
				return nil, fmt.Errorf("failed to decode recipe: %w", err)
			}
			// ran out of recipe file, all yaml documents read
			break
		}
		if err := recipe.Validate(); err != nil {
			return nil, fmt.Errorf("failed to validate recipe: %w", err)
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
		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

func renderedRecipeFilesToRecipeNames(paths []string) ([]string, error) {
	recipeNames := make([]string, len(paths))
	for i, path := range paths {
		re, err := regexp.Compile("([0-9]+)-([^.]+).yml$")
		if err != nil {
			return nil, err
		}
		parts := re.FindStringSubmatch(filepath.Base(path))
		if len(parts) != 3 {
			return nil, fmt.Errorf("Expected to match 3 parts of %s but got %d", path, len(parts))
		}
		recipeNames[i] = parts[2]
	}
	return recipeNames, nil
}
