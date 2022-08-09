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

type RenderEngine interface {
	Render(recipe *Recipe, values map[string]interface{}) (map[string][]byte, error)
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

	return nil
}

// Renders recipe templates from .Templates to .Files
func (re *Recipe) Render(engine RenderEngine) error {
	// Define the context which is available on templates
	context := map[string]interface{}{
		"Recipe":    re.Metadata,
		"Variables": re.Values,
	}

	var err error
	re.Files, err = engine.Render(re, context)
	if err != nil {
		return err
	}

	return nil
}

// Check if the recipe is in executed state (the templates has been rendered)
func (re *Recipe) IsExecuted() bool {
	return len(re.Files) > 0
}
