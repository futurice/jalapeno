package recipe

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
)

// Renders recipe templates
func (re *Recipe) Execute(values VariableValues, id uuid.UUID) (*Sauce, error) {
	if re.engine == nil {
		return nil, errors.New("render engine has not been set")
	}

	if id.IsNil() {
		return nil, errors.New("ID was nil")
	}

	sauce := NewSauce()
	sauce.Recipe = *re
	sauce.Values = values
	sauce.ID = id

	// Define the context which is available on templates
	context := map[string]interface{}{
		"ID":        sauce.ID,
		"Recipe":    re.Metadata,
		"Variables": values,
	}

	files, err := re.engine.Render(re.Templates, context)
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
