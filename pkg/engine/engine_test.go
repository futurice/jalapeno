package engine

import (
	"bytes"
	"testing"

	"github.com/futurice/jalapeno/pkg/recipe"
)

func TestRender(t *testing.T) {
	c := &recipe.Recipe{
		Metadata: recipe.Metadata{
			Name:    "test-render",
			Version: "1.2.3",
		},
		Templates: []*recipe.File{
			{Name: "templates/test1", Data: []byte("{{.var1 | title }} {{.var2 | title}}")},
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

	expect := []*recipe.File{
		{
			Name: "templates/test1",
			Data: []byte("First Second"),
		},
	}

	for i := range expect {
		if bytes.Compare(out[i].Data, expect[i].Data) != 0 {
			t.Errorf("Expected %q, got %q", expect[i].Data, out[i].Data)
		}
	}
}
