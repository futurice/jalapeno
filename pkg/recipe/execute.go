package recipe

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
)

// Execute executes the recipe and returns a sauce
func (re *Recipe) Execute(engine RenderEngine, values VariableValues, id uuid.UUID) (*Sauce, error) {
	if id.IsNil() {
		return nil, errors.New("ID was nil")
	}

	sauce := NewSauce()
	sauce.Recipe = *re
	sauce.Values = values
	sauce.ID = id

	// Define the context which is available on templates
	context := map[string]interface{}{
		"ID": sauce.ID.String(),
		"Recipe": struct{ APIVersion, Name, Version string }{
			re.APIVersion,
			re.Name,
			re.Version,
		},
		"Variables": values,
	}

	files, err := engine.Render(re.Templates, context)
	if err != nil {
		return nil, err
	}

	sauce.Files = make(map[string]File, len(files))
	idx := 0
	for filename, content := range files {
		// Skip empty files
		if len(strings.TrimSpace(string(content))) == 0 {
			continue
		}

		// Skip files starting with "_"
		if strings.HasPrefix(filename, "_") {
			continue
		}

		sum := sha256.Sum256(content)
		sauce.Files[filename] = File{Content: content, Checksum: fmt.Sprintf("sha256:%x", sum)}
		idx += 1
		if idx > len(files) {
			return nil, errors.New("files array grew during execution")
		}
	}

	if err = sauce.Validate(); err != nil {
		return nil, fmt.Errorf("sauce was not valid: %w", err)
	}

	return sauce, nil
}
