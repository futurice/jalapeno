package recipe

import (
	"fmt"
)

type Manifest struct {
	APIVersion string           `yaml:"apiVersion"`
	Recipes    []ManifestRecipe `yaml:"recipes"`
}

type ManifestRecipe struct {
	Name       string         `yaml:"name"`
	Version    string         `yaml:"version"`
	Repository string         `yaml:"repository"`
	Values     VariableValues `yaml:"values,omitempty"`
}

func NewManifest() *Manifest {
	return &Manifest{
		APIVersion: "v1",
	}
}

func (m *Manifest) Validate() error {
	if m.APIVersion != "v1" {
		return fmt.Errorf("unreconized manifest API version \"%s\"", m.APIVersion)
	}

	return nil
}

func (m *Manifest) GetRecipes() ([]*Recipe, error) {
	recipes := make([]*Recipe, len(m.Recipes))
	for i, recipe := range m.Recipes {
		re, err := LoadRecipe(recipe.Repository)
		if err != nil {
			return nil, fmt.Errorf("can not load recipe \"%s\": %w", recipe.Name, err)
		}

		recipes[i] = re
	}

	return recipes, nil
}
