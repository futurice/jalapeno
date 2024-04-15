package recipe

import (
	"fmt"

	"golang.org/x/mod/semver"
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

func NewManifest() Manifest {
	return Manifest{
		APIVersion: "v1",
	}
}

func (m *Manifest) Validate() error {
	if m.APIVersion != "v1" {
		return fmt.Errorf("unreconized manifest API version \"%s\"", m.APIVersion)
	}

	for _, r := range m.Recipes {
		if r.Name == "" {
			return fmt.Errorf("recipe name is required")
		}

		errPrefix := fmt.Sprintf("recipe '%s'", r.Name)
		if !semver.IsValid(r.Version) {
			return fmt.Errorf("%s: recipe version is not a valid semver", errPrefix)
		}

		if r.Repository == "" {
			return fmt.Errorf("%s: recipe repository is required", errPrefix)
		}

		// if {
		// 	// TODO: make sure that the repository is a valid URL
		// }

		if err := r.Values.Validate(); err != nil {
			return fmt.Errorf("%s: variable values were invalid: %w", errPrefix, err)
		}
	}

	return nil
}
