package engine

import (
	"errors"
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
		t.New(file.Name).Parse(string(file.Data))

		var buf strings.Builder
		if err := t.ExecuteTemplate(&buf, file.Name, values); err != nil {
			return map[string]string{}, errors.New("failed to execute template")
		}

		rendered[file.Name] = strings.ReplaceAll(buf.String(), "<no value>", "")
	}

	return rendered, nil
}
