package recipe

import (
	"fmt"

	"github.com/gofrs/uuid"
)

// Sauce represents a rendered recipe
type Sauce struct {
	Recipe Recipe          `yaml:",inline"`
	Values VariableValues  `yaml:"values,omitempty"`
	Files  map[string]File `yaml:"files"`

	// Random unique ID whose value is determined on first render and stays the same
	// on subsequent re-renders (upgrades) of the sauce. Can be used for example as a seed
	// for template random functions to provide same result on each template
	ID uuid.UUID `yaml:"id"`

	// CheckFrom defines the repository where updates should be checked for the recipe
	CheckFrom string `yaml:"from,omitempty"`
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
		if _, found := s.Values[variable.Name]; !(variable.Optional || variable.If != "") && !found {
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
