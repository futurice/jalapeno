package recipe

type Metadata struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	// Source is the URL to the source code of this chart
	Source      string              `json:"sources,omitempty"`
	Description string              `yaml:"description"`
	Variables   map[string]Variable `yaml:"vars,omitempty"`
}

func (re *Metadata) Validate() error {
	return nil
}
