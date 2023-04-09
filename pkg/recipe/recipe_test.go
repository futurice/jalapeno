package recipe

import (
	"html/template"
	"strings"
	"testing"

	"github.com/Masterminds/sprig"
)

type TestRenderEngine struct{}

var _ RenderEngine = TestRenderEngine{}

func (e TestRenderEngine) Render(templates map[string][]byte, values map[string]interface{}) (map[string][]byte, error) {
	t := template.New("gotpl")
	t.Funcs(sprig.TxtFuncMap())
	rendered := make(map[string][]byte)

	for name, data := range templates {
		if _, err := t.New(name).Parse(string(data)); err != nil {
			return nil, err
		}
		var buf strings.Builder
		if err := t.ExecuteTemplate(&buf, name, values); err != nil {
			return nil, err
		}
		rendered[name] = []byte(buf.String())
	}

	return rendered, nil
}

func TestRecipeRenderChecksums(t *testing.T) {
	recipe := &Recipe{
		Metadata: Metadata{
			Name:    "test",
			Version: "v1.0.0",
		},
		Variables: []Variable{
			{
				Name: "foo",
			},
		},
		Templates: map[string][]byte{
			"README.md": []byte("{{ .Variables.foo }}"),
		},
	}

	recipe.SetEngine(TestRenderEngine{})

	sauce, err := recipe.Execute(VariableValues{"foo": "bar"})
	if err != nil {
		t.Fatalf("Failed to render recipe: %s", err)
	}

	readme := sauce.Files["README.md"]
	sumWithAlgo := "sha256:fcde2b2edba56bf408601fb721fe9b5c338d10ee429ea04fae5511b68fbf8fb9"
	if readme.Checksum != sumWithAlgo {
		t.Fatalf("Expected checksum %s for content %s to match %s", readme.Content, readme.Checksum, sumWithAlgo)
	}
}
