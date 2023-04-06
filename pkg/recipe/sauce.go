package recipe

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-yaml/yaml"
)

type File struct {
	Checksum string `yaml:"checksum"` // e.g. "sha256:asdjfajdfa" w. default algo
	Content  []byte `yaml:"-"`
}

// Sauce represents a rendered recipe
type Sauce struct {
	Recipe Recipe          `yaml:",inline"`
	Values VariableValues  `yaml:"values,omitempty"`
	Files  map[string]File `yaml:"files"`
}

const (
	SauceFileName = "sauce"
	SauceDirName  = ".jalapeno"
)

func NewSauce() *Sauce {
	return &Sauce{}
}

func (s *Sauce) Validate() error {
	// TODO
	return nil
}

// Load all sauces from a project directory
func LoadSauce(projectDir string) ([]*Sauce, error) {
	var sauces []*Sauce

	sauceFile := filepath.Join(projectDir, SauceDirName, SauceFileName+YAMLExtension)
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
