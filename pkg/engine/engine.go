package engine

import (
	"errors"
	"fmt"
	"strings"
	"text/template"
)

type Engine struct {
}

func (e Engine) Render(templates map[string][]byte, values map[string]interface{}) (map[string][]byte, error) {
	t := template.New("gotpl")
	t.Funcs(funcMap())

	rendered := make(map[string][]byte)

	// Parse all templates first
	for name, data := range templates {
		_, err := t.New(name).Parse(string(data))
		if err != nil {
			// TODO: Inner error message includes prefix "template: ", which does not good when printing this error
			return nil, fmt.Errorf("failed to parse template: %w", err)
		}
	}

	// Execute each template seperately
	for name := range templates {
		var buf strings.Builder
		if err := t.ExecuteTemplate(&buf, name, values); err != nil {
			// TODO: Inner error message includes prefix "template: ", which does not good when printing this error
			return nil, fmt.Errorf("failed to execute template: %w", err)
		}

		output := buf.String()

		// If template uses variables which were undefined, gotpl will insert "<no value>"
		if strings.Contains(output, "<no value>") {
			// TODO: Find out which variable was not defined
			return nil, errors.New("some of the variables used in the template were undefined")
		}

		// TODO: Could we detect unused variables, and give warning about those?

		rendered[name] = []byte(output)
	}

	return rendered, nil
}
