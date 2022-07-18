package recipe

type Metadata struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
	URL         string `yaml:"url,omitempty"`
}

func (m *Metadata) Validate() error {
	// TODO
	return nil
}
