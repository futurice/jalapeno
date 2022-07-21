package engine

import (
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/futurice/jalapeno/pkg/recipe"
)

type Engine struct {
}

func Render(recipe *recipe.Recipe, values map[string]interface{}) ([]*recipe.File, error) {
	return new(Engine).Render(recipe, values)
}

func (e Engine) Render(r *recipe.Recipe, values map[string]interface{}) ([]*recipe.File, error) {
	t := template.New("gotpl")
	t.Funcs(funcMap())

	rendered := make([]*recipe.File, 0, len(r.Templates))

	for _, file := range r.Templates {
		_, err := t.New(file.Name).Parse(string(file.Data))
		if err != nil {
			// TODO: Inner error message includes prefix "template: ", which does not good when printing this error
			return nil, fmt.Errorf("failed to parse template: %w", err)
		}

		var buf strings.Builder
		if err := t.ExecuteTemplate(&buf, file.Name, values); err != nil {
			// TODO: Inner error message includes prefix "template: ", which does not good when printing this error
			return nil, fmt.Errorf("failed to execute template: %w", err)
		}

		output := buf.String()

		// If template uses variables which were undefined, gotpl will insert "<no value>"
		if strings.Contains(output, "<no value>") {
			// TODO: Find out which variable was not defined
			return nil, errors.New("some of the variables used in the template were undefined")
		}

		rendered = append(rendered, &recipe.File{Name: file.Name, Data: []byte(output)})
	}

	return rendered, nil
}
