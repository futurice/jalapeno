package recipe

import (
	"fmt"
)

type Recipe struct {
	Metadata  `yaml:",inline"`
	Variables []Variable        `yaml:"vars,omitempty"`
	Templates map[string][]byte `yaml:"-"`
	Tests     []Test            `yaml:"-"`
}

type RenderEngine interface {
	Render(templates map[string][]byte, values map[string]interface{}) (map[string][]byte, error)
}

func NewRecipe() *Recipe {
	return &Recipe{
		Metadata: Metadata{
			APIVersion: "v1",
		},
	}
}

func (re *Recipe) Validate() error {
	if err := re.Metadata.Validate(); err != nil {
		return err
	}

	checkDuplicates := make(map[string]bool)
	for _, v := range re.Variables {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("error on variable %s: %w", v.Name, err)
		}
		if _, exists := checkDuplicates[v.Name]; exists {
			return fmt.Errorf("variable %s has been declared multiple times", v.Name)
		}
		checkDuplicates[v.Name] = true
	}

	for _, t := range re.Tests {
		if err := t.Validate(); err != nil {
			return fmt.Errorf("error when validating recipe test case %s: %w", t.Name, err)
		}
	}

	return nil
}
