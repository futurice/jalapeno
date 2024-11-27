package engine

import (
	"errors"
	"fmt"
	"strings"
	"text/template"
)

type Engine struct {
}

func New() Engine {
	return Engine{}
}

func (e Engine) Render(templates map[string][]byte, values map[string]interface{}) (map[string][]byte, error) {
	t := template.New("gotpl")
	t.Funcs(funcMap(t))

	rendered := make(map[string][]byte)

	// Parse all templates first
	for name, data := range templates {
		_, err := t.New(name).Parse(string(data))
		if err != nil {
			// TODO: Inner error message includes prefix "template: ", which does not look good when printing the error
			return nil, fmt.Errorf("failed to parse template: %w", err)
		}
	}

	// Execute each template seperately
	for name := range templates {
		var buf strings.Builder
		if err := t.ExecuteTemplate(&buf, name, values); err != nil {
			// TODO: Inner error message includes prefix "template: ", which does not look good when printing the error
			return nil, fmt.Errorf("failed to execute template: %w", err)
		}

		output := buf.String()

		// If template uses variables which were undefined, gotpl will insert "<no value>"
		if strings.Contains(output, "<no value>") {
			// TODO: Find out which variable was not defined
			return nil, errors.New("some of the variables used in the template were undefined")
		}

		// TODO: Could we detect unused variables, and give warning about those?

		// File names can be templates to, render them here as well

		buf.Reset()

		template, err := t.New("__template_file_filename").Parse(name)
		if err != nil {
			return nil, fmt.Errorf("failed to parse file name template: %w", err)
		}

		if err := template.Execute(&buf, values); err != nil {
			return nil, fmt.Errorf("failed to execute file name template: %w", err)
		}

		filename := buf.String()
		if strings.Contains(filename, "<no value>") {
			return nil, fmt.Errorf("variable for file name %q was undefined", name)
		}

		rendered[filename] = []byte(output)
	}

	return rendered, nil
}
