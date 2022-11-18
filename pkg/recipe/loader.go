package recipe

import (
	"fmt"
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

// Load recipe which already has been rendered by name.
// The stack index of the rendered recipe gets discovered
// automatically.
func LoadRendered(path, recipeName string) (*Recipe, error) {
	rootDir, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	matches, err := filepath.Glob(filepath.Join(rootDir, RenderedRecipeDirName, fmt.Sprintf("*-%s.yml", recipeName)))
	if err != nil {
		// The only case for this should be a malformed glob pattern
		return nil, err
	} else if len(matches) != 1 {
		return nil, fmt.Errorf("Directory %s does not contain recipe %s", path, recipeName)
	}

	recipeFile := matches[0]
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

	// read rendered files
	for path, file := range recipe.Files {
		data, err := os.ReadFile(filepath.Join(rootDir, path))
		if err != nil {
			return nil, err
		}
		file.Content = data
		recipe.Files[path] = file
	}

	return recipe, nil
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

func LoadAllRendered(path string) ([]*Recipe, error) {
	matches, err := filepath.Glob(filepath.Join(path, RenderedRecipeDirName, "*-*.yml"))
	if err != nil {
		return nil, err
	}

	recipes := make([]*Recipe, len(matches))
	recipeNames, err := renderedRecipeFilesToRecipeNames(matches)
	if err != nil {
		return nil, err
	}

	for i, recipeName := range recipeNames {
		recipe, err := LoadRendered(path, recipeName)
		if err != nil {
			return nil, err
		}
		recipes[i] = recipe
	}
	return recipes, err
}
