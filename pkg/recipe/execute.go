package recipe

import (
	"errors"
	"fmt"
	"maps"
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

	// Filter out templates we might not want to render
	templates := make(map[string][]byte)
	plainFiles := make(map[string][]byte)
	for filename, file := range re.Templates {
		if strings.HasSuffix(filename, re.TemplateExtension) {
			templates[filename] = file.Content
		} else {
			plainFiles[filename] = file.Content
		}
	}

	files, err := engine.Render(templates, context)
	if err != nil {
		return nil, err
	}

	// Add the plain files
	maps.Copy(files, plainFiles)

	sauce.Files = make(map[string]File, len(re.Templates))

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

		filename = strings.TrimSuffix(filename, re.TemplateExtension)

		sauce.Files[filename] = NewFile(content)
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
