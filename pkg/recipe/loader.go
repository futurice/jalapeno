package recipe

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-yaml/yaml"
)

const (
	RecipeFileName         = "recipe.yml"
	RecipeTemplatesDirName = "templates"
)

func Load(path string) (*Recipe, error) {
	// Later on here we can add additional load mechanisms (example from URL)
	return LoadFromDir(path)
}

func LoadFromDir(path string) (*Recipe, error) {
	rootdir, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	// TODO: Check if root path was not a dir

	recipeFile := filepath.Join(rootdir, RecipeFileName)
	dat, err := os.ReadFile(recipeFile)
	if err != nil {
		return nil, err
	}

	recipe := &Recipe{}

	err = yaml.Unmarshal(dat, recipe)
	if err != nil {
		return nil, err
	}

	templates := make([]*File, 0)

	templatesDir := filepath.Join(rootdir, RecipeTemplatesDirName)

	walk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Continue walking if the path is directory
		if info.IsDir() {
			return nil
		}

		templateData, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		prefix := filepath.Join(rootdir, RecipeTemplatesDirName)
		prefix += string(filepath.Separator)
		templateName := filepath.ToSlash(strings.TrimPrefix(path, prefix))

		file := &File{
			Name: templateName,
			Data: templateData,
		}

		templates = append(templates, file)
		return nil
	}

	err = filepath.Walk(templatesDir, walk)
	if err != nil {
		return recipe, err
	}

	recipe.Templates = templates

	return recipe, nil
}
