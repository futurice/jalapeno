package recipe

import "fmt"

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

func (m *Manifest) Validate() error {
	if m.APIVersion != "v1" {
		return fmt.Errorf("unreconized manifest API version \"%s\"", m.APIVersion)
	}

	return nil
}

func (m *Manifest) GetRecipes() ([]*Recipe, error) {
	// TODO
	return nil, nil
}
