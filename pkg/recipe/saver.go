package recipe

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	defaultFileMode os.FileMode = 0700
	yamlIndent      int         = 2
)

// Save saves recipe to given destination
func (re *Recipe) Save(dest string) error {
	// TODO: Make sure recipe name is path friendly
	recipeDir := filepath.Join(dest, re.Name)

	err := os.MkdirAll(recipeDir, defaultFileMode)
	if err != nil {
		return fmt.Errorf("can not create directory %s: %v", recipeDir, err)
	}

	recipeFilepath := filepath.Join(recipeDir, RecipeFileName+YAMLExtension)
	file, err := os.Create(recipeFilepath)
	if err != nil {
		return fmt.Errorf("failed to create recipe file: %w", err)
	}

	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(yamlIndent)
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

	testRootDir := filepath.Join(dest, RecipeTestsDirName)

	err := os.MkdirAll(testRootDir, defaultFileMode)
	if err != nil {
		return fmt.Errorf("can not create recipe test directory: %w", err)
	}

	for _, test := range re.Tests {
		testDirPath := filepath.Join(testRootDir, test.Name)
		err := os.MkdirAll(filepath.Join(testRootDir, test.Name), defaultFileMode)
		if err != nil {
			return fmt.Errorf("failed to create test directory for test '%s': %w", test.Name, err)
		}

		meta, err := os.Create(filepath.Join(testDirPath, RecipeTestMetaFileName+YAMLExtension))
		if err != nil {
			return fmt.Errorf("failed to create recipe test file: %w", err)
		}
		defer meta.Close()

		encoder := yaml.NewEncoder(meta)
		encoder.SetIndent(yamlIndent)
		defer encoder.Close()

		if err := encoder.Encode(test); err != nil {
			return fmt.Errorf("failed to write recipe test to a file: %w", err)
		}

		testFileDirPath := filepath.Join(testDirPath, RecipeTestFilesDirName)
		if len(test.Files) > 0 {
			err := os.MkdirAll(testFileDirPath, defaultFileMode)
			if err != nil {
				return fmt.Errorf("failed to create test file directory for test '%s': %w", test.Name, err)
			}

			err = saveFileMap(test.Files, testFileDirPath)
			if err != nil {
				return fmt.Errorf("failed to save template files: %w", err)
			}
		}
	}

	return nil
}

func (re *Recipe) saveTemplates(dest string) error {
	templateDir := filepath.Join(dest, RecipeTemplatesDirName)
	err := os.MkdirAll(templateDir, defaultFileMode)
	if err != nil {
		return fmt.Errorf("can not save templates to the directory: %w", err)
	}

	err = saveFileMap(re.Templates, templateDir)
	if err != nil {
		return fmt.Errorf("failed to save template files: %w", err)
	}

	return nil
}

// Save saves sauce to given destination
func (s *Sauce) Save(dest string) error {
	// load all sauces from target dir, because we will either replace
	// a previous sauce, or create a new file
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
		return fmt.Errorf("failed to create sauce dir: %w", err)
	}
	file, err := os.Create(filepath.Join(dest, SauceDirName, SaucesFileName+YAMLExtension))
	if err != nil {
		return fmt.Errorf("failed to create sauce file: %w", err)
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(yamlIndent)
	defer encoder.Close()

	for _, sauce := range sauces {
		if err := encoder.Encode(sauce); err != nil {
			return fmt.Errorf("failed to write sauces: %w", err)
		}
	}

	fileMap := make(map[string][]byte)
	for filename, file := range s.Files {
		fileMap[filename] = file.Content
	}

	err = saveFileMap(fileMap, dest)
	if err != nil {
		return fmt.Errorf("failed to save sauce files: %w", err)
	}

	return nil
}

func saveFileMap(files map[string][]byte, dest string) error {
	if len(files) == 0 {
		return nil
	}

	if _, err := os.Stat(dest); os.IsNotExist(err) {
		return errors.New("destination path does not exist")
	}

	for path, file := range files {
		destPath := filepath.Join(dest, path)

		// Create file's parent directories (if not already exist)
		err := os.MkdirAll(filepath.Dir(destPath), 0700)
		if err != nil {
			return err
		}

		// Create the file
		f, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer f.Close()

		// Write the data to the file
		_, err = f.Write(file)
		if err != nil {
			return err
		}

		err = f.Sync()
		if err != nil {
			return err
		}
	}
	return nil
}
