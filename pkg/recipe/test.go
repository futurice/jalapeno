package recipe

type Test struct {
	Name   string            `yaml:"name,omitempty"`
	Values map[string]string `yaml:"values,omitempty"`
	Files  map[string]string `yaml:"files"`
}

func (t *Test) Validate() error {
	// TODO
	return nil
}
