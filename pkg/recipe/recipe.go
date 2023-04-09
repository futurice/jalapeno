package recipe

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/futurice/jalapeno/pkg/engine"
)

type Recipe struct {
	Metadata  `yaml:",inline"`
	Variables []Variable        `yaml:"vars,omitempty"`
	Templates map[string][]byte `yaml:"-"`
	Tests     []Test            `yaml:"-"`
	engine    RenderEngine
}

type RenderEngine interface {
	Render(templates map[string][]byte, values map[string]interface{}) (map[string][]byte, error)
}

func NewRecipe() *Recipe {
	return &Recipe{
		Metadata: Metadata{
			APIVersion: "v1",
		},
		engine: engine.Engine{},
	}
}

func (re *Recipe) Validate() error {
	if err := re.Metadata.Validate(); err != nil {
		return err
	}

	checkDuplicates := make(map[string]bool)
	for _, v := range re.Variables {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("error on variable %s: %w", v.Name, err)
		}
		if _, exists := checkDuplicates[v.Name]; exists {
			return fmt.Errorf("variable %s has been declared multiple times", v.Name)
		}
		checkDuplicates[v.Name] = true
	}

	for _, t := range re.Tests {
		if err := t.Validate(); err != nil {
			return fmt.Errorf("error when validating recipe test case %s: %w", t.Name, err)
		}
	}

	return nil
}

func (re *Recipe) SetEngine(e RenderEngine) {
	re.engine = e
}

// Renders recipe templates from .Templates to .Files
func (re *Recipe) Execute(values VariableValues) (*Sauce, error) {
	if re.engine == nil {
		return nil, errors.New("render engine has not been set")
	}

	// Define the context which is available on templates
	context := map[string]interface{}{
		"Recipe":    re.Metadata,
		"Variables": values,
	}

	var err error
	files, err := re.engine.Render(re.Templates, context)
	if err != nil {
		return nil, err
	}

	sauce := &Sauce{
		Recipe: *re,
		Values: values,
	}

	sauce.Files = make(map[string]File, len(files))
	idx := 0
	for filename, content := range files {
		sum := sha256.Sum256(content)
		sauce.Files[filename] = File{Content: content, Checksum: fmt.Sprintf("sha256:%x", sum)}
		idx += 1
		if idx > len(files) {
			return nil, errors.New("files array grew during execution")
		}
	}

	return sauce, nil
}

type RecipeConflict struct {
	Path           string
	Sha256Sum      string
	OtherSha256Sum string
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
