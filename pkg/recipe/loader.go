package recipe

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	YAMLExtension          = ".yml"
	RecipeFileName         = "recipe"
	RecipeTemplatesDirName = "templates"
	RecipeTestsDirName     = "tests"
	RecipeTestMetaFileName = "test"
	RecipeTestFilesDirName = "files"
	IgnoreFileName         = ".jalapenoignore"
)

var (
	ErrSauceNotFound = errors.New("sauce not found")
)

// LoadRecipe reads a recipe from a given path
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

func loadTemplates(recipePath string) (map[string]File, error) {
	templates := make(map[string]File)

	walk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Continue walking if the path is directory
		if info.IsDir() {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Create a filepath related to the root of the directory
		prefix := fmt.Sprintf("%s%c", recipePath, filepath.Separator)
		name := filepath.ToSlash(strings.TrimPrefix(path, prefix))

		templates[name] = NewFile(content)
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

	testDirs, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, dir := range testDirs {
		if !dir.IsDir() {
			continue
		}

		test := Test{}
		testDirPath := filepath.Join(path, dir.Name())
		contents, err := os.ReadFile(filepath.Join(testDirPath, RecipeTestMetaFileName+YAMLExtension))
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(contents, &test)
		if err != nil {
			return nil, err
		}

		test.Name = dir.Name()
		test.Files = make(map[string]File)
		testFileDirPath := filepath.Join(testDirPath, RecipeTestFilesDirName)

		walk := func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			// Create a filepath related to the root of the test file directory
			prefix := fmt.Sprintf("%s%c", testFileDirPath, filepath.Separator)
			trimmedPath := filepath.ToSlash(strings.TrimPrefix(path, prefix))

			test.Files[trimmedPath] = NewFile(content)

			return nil
		}

		err = filepath.Walk(testFileDirPath, walk)
		if err != nil {
			return nil, fmt.Errorf("error when loading test files: %w", err)
		}

		tests = append(tests, test)
	}

	return tests, nil
}

func NewFile(content []byte) File {
	return File{
		Content:  content,
		Checksum: fmt.Sprintf("sha256:%x", sha256.Sum256(content)),
	}
}
