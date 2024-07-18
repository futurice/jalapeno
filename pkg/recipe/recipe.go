package recipe

import (
	"fmt"
)

type Recipe struct {
	Metadata  `yaml:",inline"`
	Variables []Variable      `yaml:"vars,omitempty"`
	Templates map[string]File `yaml:"-"`
	Tests     []Test          `yaml:"-"`
}

func NewRecipe() Recipe {
	return Recipe{
		Metadata: Metadata{
			APIVersion: "v1",
		},
	}
}

func (re *Recipe) Validate() error {
	if err := re.Metadata.Validate(); err != nil {
		return err
	}

	varDuplicateCheck := make(map[string]struct{})
	for _, v := range re.Variables {
		if _, exists := varDuplicateCheck[v.Name]; exists {
			return fmt.Errorf("variable %s has been declared multiple times", v.Name)
		}
		if err := v.Validate(); err != nil {
			return fmt.Errorf("error on variable %s: %w", v.Name, err)
		}
		varDuplicateCheck[v.Name] = struct{}{}
	}

	testDuplicateCheck := make(map[string]struct{})
	for _, t := range re.Tests {
		if _, exists := testDuplicateCheck[t.Name]; exists {
			return fmt.Errorf("test case %s has been declared multiple times", t.Name)
		}

		if err := t.Validate(); err != nil {
			return fmt.Errorf("error when validating recipe test case %s: %w", t.Name, err)
		}
		testDuplicateCheck[t.Name] = struct{}{}

		for vName := range t.Values {
			found := false
			for _, v := range re.Variables {
				if v.Name == vName {
					found = true
					break
				}
			}

			if !found {
				return fmt.Errorf("test case %s references an unknown variable %s", t.Name, vName)
			}
		}
	}

	return nil
}
