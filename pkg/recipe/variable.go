package recipe

type Variable struct {
	Placeholder string            `yaml:"placeholder,omitempty"`
	Options     []string          `yaml:"options,omitempty"`
	Template    string            `yaml:"template,omitempty"`
	Validator   VariableValidator `yaml:"validator,omitempty"`
}

type VariableValidator struct {
	Type   string `yaml:"type,omitempty"`
	RegExp string `yaml:"regexp,omitempty"`
}

type VariableValues map[string]string
