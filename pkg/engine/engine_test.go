package engine

import (
	"testing"

	"github.com/futurice/jalapeno/pkg/recipe"
)

func TestRender(t *testing.T) {
	c := &recipe.Recipe{
		Metadata: &recipe.Metadata{
			Name:    "test-render",
			Version: "1.2.3",
		},
		Templates: []*recipe.File{
			{Name: "templates/test1", Data: []byte("{{.var1 | title }} {{.var2 | title}}")},
			{Name: "templates/test2", Data: []byte("{{.noValue}}")},
		},
	}

	vals := map[string]interface{}{
		"var1": "first",
		"var2": "second",
	}

	out, err := Render(c, vals)

	if err != nil {
		t.Errorf("Failed to render templates: %s", err)
	}

	expect := map[string]string{
		"templates/test1": "First Second",
		"templates/test2": "",
	}

	for name, data := range expect {
		if out[name] != data {
			t.Errorf("Expected %q, got %q", data, out[name])
		}
	}
}
