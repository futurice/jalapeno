package recipe

import (
	"testing"

	"github.com/gofrs/uuid"
)

func TestRenderInitHelp(t *testing.T) {
	scenarios := []struct {
		name           string
		help           string
		values         VariableValues
		expectedOutput string
		expectingErr   bool
	}{
		{
			"conditional text",
			"{{ if .Variables.FOO }}Foo is true{{ else }}Foo is false{{ end }}",
			VariableValues{
				"FOO": true,
			},
			"Foo is true",
			false,
		},
		{
			"invalid template",
			"{{ if .Variables.NOT_FOUND }}",
			VariableValues{
				"FOO": true,
			},
			"",
			true,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			re := NewRecipe()
			re.Name = "test"
			re.Version = "v0.0.0"
			re.InitHelp = scenario.help

			sauce := NewSauce()
			sauce.Recipe = *re
			sauce.ID = uuid.Must(uuid.NewV4())
			sauce.Values = scenario.values

			if help, err := sauce.RenderInitHelp(); err != nil {
				if !scenario.expectingErr {
					t.Fatalf("Got error when not expected: %s", err)
				}
			} else if scenario.expectingErr {
				t.Fatal("Was expecting error, did not get any")

			} else if help != scenario.expectedOutput {
				t.Fatalf("Expected output '%s', got '%s'", scenario.expectedOutput, help)
			}
		})
	}
}
