package recipe

import (
	"html/template"
	"strings"
	"testing"

	"github.com/Masterminds/sprig"
)

type TestRenderEngine struct{}

var _ RenderEngine = TestRenderEngine{}

func (e TestRenderEngine) Render(r *Recipe, values map[string]interface{}) (map[string][]byte, error) {
	t := template.New("gotpl")
	t.Funcs(sprig.TxtFuncMap())
	rendered := make(map[string][]byte)

	for name, data := range r.Templates {
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
		Values: VariableValues{"foo": "bar"},
		Templates: map[string][]byte{
			"README.md": []byte("{{ .Variables.foo }}"),
		},
	}

	if err := recipe.Render(TestRenderEngine{}); err != nil {
		t.Error("Failed to render recipe", err)
	}

	readme := recipe.Files["README.md"]
	sumWithAlgo := "sha256:fcde2b2edba56bf408601fb721fe9b5c338d10ee429ea04fae5511b68fbf8fb9"
	if readme.Checksum != sumWithAlgo {
		t.Errorf("Expected checksum %s for content %s to match %s", readme.Content, readme.Checksum, sumWithAlgo)
	}
}
