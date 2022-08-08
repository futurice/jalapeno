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
		Templates: map[string][]byte{
			"templates/test1": []byte("{{.var1 | title }} {{.var2 | title}}"),
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

	expect := map[string][]byte{
		"templates/test1": []byte("First Second"),
	}

	for name := range expect {
		if bytes.Equal(out[name], expect[name]) {
			t.Errorf("Expected %q, got %q", expect[name], out[name])
		}
	}
}
