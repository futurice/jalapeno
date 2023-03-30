package engine

import (
	"bytes"
	"testing"
)

func TestRender(t *testing.T) {
	templates := map[string][]byte{
		"templates/test1": []byte("{{.var1 | title }} {{.var2 | title}}"),
	}

	vals := map[string]interface{}{
		"var1": "first",
		"var2": "second",
	}

	out, err := new(Engine).Render(templates, vals)

	if err != nil {
		t.Fatalf("Failed to render templates: %s", err)
	}

	expect := map[string][]byte{
		"templates/test1": []byte("First Second"),
	}

	for name := range expect {
		if !bytes.Equal(out[name], expect[name]) {
			t.Fatalf("Expected %q, got %q", expect[name], out[name])
		}
	}
}
