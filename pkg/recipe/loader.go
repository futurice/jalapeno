package recipe

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofrs/uuid"
	"gopkg.in/yaml.v3"
)

const (
	YAMLExtension    = ".yml"
	MetadataFileName = "recipe"
	TemplatesDirName = "templates"
	TestsDirName     = "tests"
	TestMetaFileName = "test"
	TestFilesDirName = "files"
	IgnoreFileName   = ".jalapenoignore"
	ManifestFileName = "manifest"
)

var (
	ErrSauceNotFound  = errors.New("sauce not found")
	ErrAmbiguousSauce = errors.New("multiple sauces found with same recipe")
)

// LoadRecipe reads a recipe from a given path
func LoadRecipe(path string) (*Recipe, error) {
	rootDir, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	recipeFile := filepath.Join(rootDir, MetadataFileName+YAMLExtension)
	dat, err := os.ReadFile(recipeFile)
	if err != nil {
		return nil, err
	}

	recipe := NewRecipe()
	err = yaml.Unmarshal(dat, &recipe)
	if err != nil {
		return nil, fmt.Errorf("error when reading %s file: %w", MetadataFileName+YAMLExtension, err)
	}

	recipe.Templates, err = loadTemplates(filepath.Join(rootDir, TemplatesDirName))
	if err != nil {
		return nil, fmt.Errorf("error when loading recipe templates: %w", err)
	}

	recipe.Tests, err = loadTests(filepath.Join(rootDir, TestsDirName))
	if err != nil {
		return nil, fmt.Errorf("error when loading recipe tests: %w", err)
	}

	if err := recipe.Validate(); err != nil {
		return nil, fmt.Errorf("loaded recipe was invalid: %w", err)
	}

	return &recipe, nil
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
		contents, err := os.ReadFile(filepath.Join(testDirPath, TestMetaFileName+YAMLExtension))
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(contents, &test)
		if err != nil {
			return nil, err
		}

		test.Name = dir.Name()
		test.Files = make(map[string]File)
		testFileDirPath := filepath.Join(testDirPath, TestFilesDirName)

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

// Load all sauces from a project directory. Returns empty slice if the project directory did not contain any sayces
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
	sauceData, err := os.ReadFile(sauceFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read recipe file: %w", err)
	}

	decoder := yaml.NewDecoder(bytes.NewReader(sauceData))
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
			data, err := os.ReadFile(filepath.Join(projectDir, filepath.Clean(sauce.Subpath), path))
			if err != nil {
				// The file have been removed by the user after the sauce was created
				if errors.Is(err, os.ErrNotExist) {
					delete(sauce.Files, path)
					continue
				} else {
					return nil, fmt.Errorf("failed to read rendered file: %w", err)
				}
			}
			// Note that we use the file checksum from the sauce file,
			// but read contents from the files. This means that if the checksum
			// does not match the data, the file has been modified outside of Jalapeno
			file.Content = data
			sauce.Files[path] = file
		}

		if err := sauce.Validate(); err != nil {
			return nil, fmt.Errorf("failed to validate sauce: %w", err)
		}

		sauces = append(sauces, sauce)
	}

	return sauces, nil
}

func LoadSauceByRecipe(projectDir, recipeName string) (*Sauce, error) {
	sauces, err := LoadSauces(projectDir)
	if err != nil {
		return nil, err
	}

	var found *Sauce
	for _, s := range sauces {
		if s.Recipe.Name == recipeName {
			if found != nil {
				return nil, fmt.Errorf("%w '%s'", ErrAmbiguousSauce, recipeName)
			}
			found = s
		}
	}

	if found != nil {
		return found, nil
	}

	return nil, ErrSauceNotFound
}

func LoadSauceByID(projectDir string, id uuid.UUID) (*Sauce, error) {
	sauces, err := LoadSauces(projectDir)
	if err != nil {
		return nil, err
	}

	for _, s := range sauces {
		if s.ID == id {
			return s, nil
		}
	}

	return nil, ErrSauceNotFound
}

func LoadManifest(path string) (*Manifest, error) {
	dat, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	manifest := &Manifest{}
	err = yaml.Unmarshal(dat, manifest)
	if err != nil {
		return nil, err
	}

	if err := manifest.Validate(); err != nil {
		return nil, err
	}

	return manifest, nil
}

// SaveSauces saves the given sauces to the project directory
func SaveSauces(projectDir string, sauces []*Sauce) error {
	if err := os.MkdirAll(filepath.Join(projectDir, SauceDirName), 0755); err != nil {
		return fmt.Errorf("failed to create sauce directory: %w", err)
	}

	sauceFile := filepath.Join(projectDir, SauceDirName, SaucesFileName+YAMLExtension)
	f, err := os.Create(sauceFile)
	if err != nil {
		return fmt.Errorf("failed to create sauce file: %w", err)
	}
	defer f.Close()

	encoder := yaml.NewEncoder(f)
	defer encoder.Close()

	for _, sauce := range sauces {
		if err := encoder.Encode(sauce); err != nil {
			return fmt.Errorf("failed to encode sauce: %w", err)
		}
	}

	return nil
}
