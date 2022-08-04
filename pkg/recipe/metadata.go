package recipe

import (
	"errors"
	"fmt"
	"net/url"

	"golang.org/x/mod/semver"
)

type Metadata struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
	URL         string `yaml:"url,omitempty"`
	InitHelp    string `yaml:"initHelp,omitempty"`
}

func (m *Metadata) Validate() error {
	if !semver.IsValid(m.Version) {
		return errors.New("version is not a valid semver")
	}

	if m.URL != "" {
		if _, err := url.ParseRequestURI(m.URL); err != nil {
			return fmt.Errorf("url is invalid: %w", err)
		}
	}

	return nil
}
