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
				"templates/test1": []byte("{{.var1 | title }} {{.var2 | title}}"),
			},
			map[string]interface{}{
				"var1": "first",
				"var2": "second",
			},
			map[string][]byte{
				"templates/test1": []byte("First Second"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			out, err := new(Engine).Render(tc.Templates, tc.Values)
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
