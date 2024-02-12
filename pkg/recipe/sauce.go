package recipe

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/fatih/structs"
	"github.com/gofrs/uuid"
)

// Sauce represents a rendered recipe
type Sauce struct {
	// Version of the sauce API schema. Currently should have value "v1"
	APIVersion string `yaml:"apiVersion"`

	// The recipe which was used to render the sauce
	Recipe Recipe `yaml:"recipe"`

	// Values which was used to execute the recipe
	Values VariableValues `yaml:"values,omitempty"`

	// Files genereated from the recipe
	Files map[string]File `yaml:"files"`

	// Random unique ID whose value is determined on first render and stays the same
	// on subsequent re-renders (upgrades) of the sauce. Can be used for example as a seed
	// for template random functions to provide same result on each template
	ID uuid.UUID `yaml:"id"`

	// CheckFrom defines the repository where updates should be checked for the recipe
	CheckFrom string `yaml:"from,omitempty"`
}

type RecipeConflict struct {
	Path           string
	Sha256Sum      string
	OtherSha256Sum string
}

const (
	SaucesFileName = "sauces"

	// The directory name which contains all Jalapeno related files
	// in the project directory
	SauceDirName = ".jalapeno"
)

func NewSauce() *Sauce {
	return &Sauce{
		APIVersion: "v1",
	}
}

func (s *Sauce) Validate() error {
	if s.APIVersion != "v1" {
		return fmt.Errorf("unreconized sauce API version \"%s\"", s.APIVersion)
	}

	if s.ID.IsNil() {
		return fmt.Errorf("sauce ID was empty")
	}

	if s.CheckFrom != "" && !strings.HasPrefix(s.CheckFrom, "oci://") {
		return fmt.Errorf("currently recipe updates can only be checked from OCI repositories, got: %s", s.CheckFrom)
	}

	if err := s.Recipe.Validate(); err != nil {
		return fmt.Errorf("sauce recipe was invalid: %w", err)
	}

	for _, variable := range s.Recipe.Variables {
		if _, found := s.Values[variable.Name]; !(variable.Optional || variable.If != "") && !found {
			return fmt.Errorf("sauce did not have value for required variable '%s'", variable.Name)
		}
	}
	return nil
}

// Check if the recipe conflicts with another recipe. Recipes conflict if they touch the same files.
func (s *Sauce) Conflicts(other *Sauce) []RecipeConflict {
	var conflicts []RecipeConflict
	for path, file := range s.Files {
		if otherFile, exists := other.Files[path]; exists {
			conflicts = append(
				conflicts,
				RecipeConflict{
					Path:           path,
					Sha256Sum:      file.Checksum,
					OtherSha256Sum: otherFile.Checksum,
				})
		}
	}
	return conflicts
}

func (s *Sauce) CreateTemplateContext() (map[string]interface{}, error) {
	if err := s.Validate(); err != nil {
		return nil, err
	}

	mappedValues := make(VariableValues)
	for name, value := range s.Values {
		switch value := value.(type) {
		// Map table to more convenient format
		case TableValue:
			mappedValues[name] = value.ToMapSlice()
		default:
			mappedValues[name] = value
		}
	}

	return structs.Map(TemplateContext{
		ID: s.ID.String(),
		Recipe: struct{ APIVersion, Name, Version, Source string }{
			s.Recipe.APIVersion,
			s.Recipe.Name,
			s.Recipe.Version,
			s.Recipe.Source,
		},
		Variables: mappedValues,
	}), nil
}

func (s *Sauce) RenderInitHelp() (string, error) {
	context, err := s.CreateTemplateContext()
	if err != nil {
		return "", err
	}

	t, err := template.New("initHelp").Parse(s.Recipe.InitHelp)
	if err != nil {
		return "", fmt.Errorf("failed to parse initHelp template: %w", err)
	}

	var buf strings.Builder
	if err := t.Execute(&buf, context); err != nil {
		return "", fmt.Errorf("failed to render initHelp template: %w", err)
	}

	return buf.String(), nil
}
