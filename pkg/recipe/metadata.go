package recipe

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/mod/semver"
)

var nameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]*$`)

type Metadata struct {
	// Version of the recipe metadata API schema. Currently should have value "v1"
	APIVersion string `yaml:"apiVersion"`

	// Name of the recipe
	Name string `yaml:"name"`

	// Version of the recipe. Must be valid [semver](https://semver.org/)
	Version string `yaml:"version"`

	// Description of what the recipe does
	Description string `yaml:"description"`

	// URL to source code for this recipe
	Source string `yaml:"source,omitempty"`

	// A message which will be showed to an user after a succesful recipe execution.
	// Can be used to guide the user what should be done next in the project directory.
	InitHelp string `yaml:"initHelp,omitempty"`

	// Glob patterns for ignoring generated files from future recipe upgrades. Ignored
	// files will not be regenerated even if their templates change in future versions
	// of the recipe.
	IgnorePatterns []string `yaml:"ignorePatterns,omitempty"`

	// File extension of files in "templates" directory which should be templated.
	// Files not matched by this extension will be copied as-is.
	// If left empty (the default), all files will be templated.
	TemplateExtension string `yaml:"templateExtension,omitempty"`
}

func (m *Metadata) Validate() error {
	// Currently we support only apiVersion v1
	if m.APIVersion != "v1" {
		return fmt.Errorf("unreconized metadata API version \"%s\"", m.APIVersion)
	}

	if m.Name == "" {
		return fmt.Errorf("recipe name can not be empty")
	}

	if !nameRegex.MatchString(m.Name) {
		return fmt.Errorf("recipe name can only contain letters, numbers, dashes and underscores")
	}

	if !semver.IsValid(m.Version) {
		return fmt.Errorf("version \"%s\" is not a valid semver", m.Version)
	}

	if m.TemplateExtension != "" && !strings.HasPrefix(m.TemplateExtension, ".") {
		return fmt.Errorf("template extension must start with a dot")
	}

	if m.Source != "" {
		if _, err := url.ParseRequestURI(m.Source); err != nil {
			return fmt.Errorf("source url is invalid: %w", err)
		}
	}

	return nil
}
