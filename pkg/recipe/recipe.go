package recipe

import "fmt"

type Recipe struct {
	Metadata  `yaml:",inline"`
	Variables map[string]Variable `yaml:"vars,omitempty"`
	Values    *VariableValues     `yaml:"values,omitempty"`
	Templates []*File             `yaml:"-"`
}

func (re *Recipe) Validate() error {
	if err := re.Metadata.Validate(); err != nil {
		return err
	}

	for name, variable := range re.Variables {
		if err := variable.Validate(); err != nil {
			return fmt.Errorf("error on variable %s: %w", name, err)
		}
	}

	return nil
}
