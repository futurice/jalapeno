package engine

import (
	"bytes"
	"testing"
)

func TestRender(t *testing.T) {
	testCases := []struct {
		Name           string
		Templates      map[string][]byte
		Values         map[string]interface{}
		ExpectedOutput map[string][]byte
	}{
		{
			"values_and_functions",
			map[string][]byte{
				"templates/test1":     []byte("{{.var1 | title }} {{.var2 | title}}"),
				"templates/{{.var1}}": []byte("{{.var1}}"),
			},
			map[string]interface{}{
				"var1": "first",
				"var2": "second",
			},
			map[string][]byte{
				"templates/test1": []byte("First Second"),
				"templates/first": []byte("first"),
			},
		},
		{
			"macros",
			map[string][]byte{
				"templates/main":    []byte("{{ template \"helper1\" }} {{ template \"helper2\" }} {{ include \"helper3\" . | upper }}"),
				"templates/helper1": []byte("{{ define \"helper1\" }}ONE{{ end }}"),
				"templates/helper2": []byte("{{ define \"helper2\" }}TWO{{ end }}"),
				"templates/helper3": []byte("{{ define \"helper3\" }}three{{ end }}"),
			},
			map[string]interface{}{},
			map[string][]byte{
				"templates/main": []byte("ONE TWO THREE"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			out, err := New().Render(tc.Templates, tc.Values)
			if err != nil {
				t.Fatalf("Failed to render templates: %s", err)
			}
			for name := range tc.ExpectedOutput {
				if !bytes.Equal(out[name], tc.ExpectedOutput[name]) {
					t.Fatalf("Expected %q, got %q", tc.ExpectedOutput[name], out[name])
				}
			}
		})
	}
}
