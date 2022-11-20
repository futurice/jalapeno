package recipe

import (
	"crypto/sha256"
	"fmt"
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
		t.New(name).Parse(string(data))
		var buf strings.Builder
		t.ExecuteTemplate(&buf, name, values)
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
			"README.md": []byte("{{ foo }}"),
		},
	}

	if err := recipe.Render(TestRenderEngine{}); err != nil {
		t.Error("Failed to render recipe", err)
	}

	readme := recipe.Files["README.md"]
	sum := sha256.Sum256(readme.Content)
	sumWithAlgo := fmt.Sprintf("sha256:%x", sum)
	if sumWithAlgo != readme.Checksum {
		t.Errorf("Expected file checksum %s to match %s", readme.Checksum, sumWithAlgo)
	}
}
