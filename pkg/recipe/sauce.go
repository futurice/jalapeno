package recipe

import (
	"fmt"
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
