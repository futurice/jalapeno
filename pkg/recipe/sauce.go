package recipe

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-yaml/yaml"
	"github.com/gofrs/uuid"
)

// Sauce represents a rendered recipe
type Sauce struct {
	Recipe Recipe          `yaml:",inline"`
	Values VariableValues  `yaml:"values,omitempty"`
	Files  map[string]File `yaml:"files"`

	// Random ID whose value is determined on first render and stays the same
	// on subsequent re-renders (upgrades) of the sauce. Can be used for example as a seed
	// for template random functions to provide same result on each template
	Anchor uuid.UUID `yaml:"anchor"`
}

type File struct {
	Checksum string `yaml:"checksum"` // e.g. "sha256:xxxxxxxxx" w. default algo
	Content  []byte `yaml:"-"`
}

type RecipeConflict struct {
	Path           string
	Sha256Sum      string
	OtherSha256Sum string
}

const (
	SaucesFileName = "sauces"

	// The directory name which contains all Jalapeno related files
	// in the project directory
	SauceDirName = ".jalapeno"
)

func NewSauce() *Sauce {
	return &Sauce{}
}

func (s *Sauce) Validate() error {
	if err := s.Recipe.Validate(); err != nil {
		return fmt.Errorf("sauce recipe was invalid: %w", err)
	}

	for _, variable := range s.Recipe.Variables {
		if _, found := s.Values[variable.Name]; !variable.Optional && !found {
			return fmt.Errorf("sauce did not have value for required variable '%s'", variable.Name)
		}
	}
	return nil
}

// Check if the recipe conflicts with another recipe. Recipes conflict if they touch the same files.
func (s *Sauce) Conflicts(other *Sauce) []RecipeConflict {
	var conflicts []RecipeConflict
	for path, file := range s.Files {
		if otherFile, exists := other.Files[path]; exists {
			conflicts = append(
				conflicts,
				RecipeConflict{
					Path:           path,
					Sha256Sum:      file.Checksum,
					OtherSha256Sum: otherFile.Checksum,
				})
		}
	}
	return conflicts
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

func LoadSauce(projectDir, recipeName string) (*Sauce, error) {
	sauces, err := LoadSauces(projectDir)
	if err != nil {
		return nil, err
	}

	for _, s := range sauces {
		if s.Recipe.Name == recipeName {
			return s, nil
		}
	}

	return nil, ErrSauceNotFound
}
