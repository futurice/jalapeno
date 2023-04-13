package recipe

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
)

type ExecuteOptions struct {
	UseStaticAnchor bool
}

// Renders recipe templates
func (re *Recipe) Execute(values VariableValues, opts ExecuteOptions) (*Sauce, error) {
	if re.engine == nil {
		return nil, errors.New("render engine has not been set")
	}

	sauce := NewSauce()
	sauce.Recipe = *re
	sauce.Values = values

	// Static anchor is used when running recipe tests
	if opts.UseStaticAnchor {
		// Not necessarily have to be nil UUID, can be any hardcoded UUID
		sauce.Anchor = uuid.Nil
	}

	// Define the context which is available on templates
	context := map[string]interface{}{
		"Anchor":    sauce.Anchor,
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
