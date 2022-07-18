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

func Render(recipe *recipe.Recipe, values map[string]interface{}) (map[string]string, error) {
	return new(Engine).Render(recipe, values)
}

func (e Engine) Render(recipe *recipe.Recipe, values map[string]interface{}) (map[string]string, error) {
	t := template.New("gotpl")
	t.Funcs(funcMap())

	rendered := make(map[string]string, 1)

	for _, file := range recipe.Templates {
		_, err := t.New(file.Name).Parse(string(file.Data))
		if err != nil {
			// TODO: Inner error message includes prefix "template: ", which does not good when printing this error
			return map[string]string{}, fmt.Errorf("failed to parse template: %w", err)
		}

		var buf strings.Builder
		if err := t.ExecuteTemplate(&buf, file.Name, values); err != nil {
			// TODO: Inner error message includes prefix "template: ", which does not good when printing this error
			return map[string]string{}, fmt.Errorf("failed to execute template: %w", err)
		}

		output := buf.String()

		// If template uses variables which were undefined, gotpl will insert "<no value>"
		if strings.Contains(output, "<no value>") {
			// TODO: Find out which variable was not defined
			return map[string]string{}, errors.New("some of the variables used in the template were undefined")
		}

		rendered[file.Name] = output
	}

	return rendered, nil
}
