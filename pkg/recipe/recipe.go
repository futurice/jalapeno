package recipe

import (
	"fmt"
)

type Recipe struct {
	Metadata  `yaml:",inline"`
	Variables []Variable        `yaml:"vars,omitempty"`
	Values    VariableValues    `yaml:"values,omitempty"`
	Templates map[string][]byte `yaml:"-"`
	Files     map[string][]byte `yaml:"-"`
}

func (re *Recipe) Validate() error {
	if err := re.Metadata.Validate(); err != nil {
		return err
	}

	for _, variable := range re.Variables {
		if err := variable.Validate(); err != nil {
			return fmt.Errorf("error on variable %s: %w", variable.Name, err)
		}
	}

	return nil
}

type RenderEngine interface {
	Render(recipe *Recipe, values map[string]interface{}) (map[string][]byte, error)
}

func (re *Recipe) Render(engine RenderEngine) error {
	var err error
	// Define the context which is available on templates
	context := map[string]interface{}{
		"Recipe":    re.Metadata,
		"Variables": re.Values,
	}

	re.Files, err = engine.Render(re, context)
	if err != nil {
		return err
	}

	return nil
}
