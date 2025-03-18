package recipe

import (
	"html/template"
	"strings"
	"testing"

	"github.com/Masterminds/sprig/v3"
	"github.com/gofrs/uuid"
)

type TestRenderEngine struct{}

var _ RenderEngine = TestRenderEngine{}

func (e TestRenderEngine) Render(templates map[string][]byte, values map[string]any) (map[string][]byte, error) {
	rendered := make(map[string][]byte)

	for name, data := range templates {
		t := template.New("gotpl")
		t.Funcs(sprig.TxtFuncMap())

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
	recipe := NewRecipe()
	recipe.Metadata.Name = "test"
	recipe.Metadata.Version = "v1.0.0"
	recipe.Variables = []Variable{{Name: "foo"}}
	recipe.Templates = map[string]File{
		"README.md": NewFile([]byte("{{ .Variables.foo }}")),
	}

	sauce, err := recipe.Execute(TestRenderEngine{}, VariableValues{"foo": "bar"}, uuid.Must(uuid.NewV4()))
	if err != nil {
		t.Fatalf("Failed to render recipe: %s", err)
	}

	readme := sauce.Files["README.md"]
	sumWithAlgo := "sha256:fcde2b2edba56bf408601fb721fe9b5c338d10ee429ea04fae5511b68fbf8fb9"
	if readme.Checksum != sumWithAlgo {
		t.Fatalf("Expected checksum %s for content %s to match %s", readme.Content, readme.Checksum, sumWithAlgo)
	}
}

func TestRecipeRenderID(t *testing.T) {
	recipe := NewRecipe()
	recipe.Metadata.Name = "test"
	recipe.Metadata.Version = "v1.0.0"

	sauce, err := recipe.Execute(TestRenderEngine{}, nil, uuid.Must(uuid.NewV4()))
	if err != nil {
		t.Fatalf("Failed to render recipe: %s", err)
	}

	if sauce.ID.IsNil() {
		t.Fatal("Sauce ID was empty")
	}
}

func TestRecipeRenderIDReuse(t *testing.T) {
	recipe := NewRecipe()
	recipe.Metadata.Name = "test"
	recipe.Metadata.Version = "v1.0.0"

	sauce1, err := recipe.Execute(TestRenderEngine{}, nil, TestID)
	if err != nil {
		t.Fatalf("Failed to render first recipe: %s", err)
	}

	sauce2, err := recipe.Execute(TestRenderEngine{}, nil, TestID)
	if err != nil {
		t.Fatalf("Failed to render second recipe: %s", err)
	}

	if sauce1.ID != sauce2.ID {
		t.Fatal("IDs were not same when used static ID on both executions")
	}
}

func TestRecipeRenderEmptyFiles(t *testing.T) {
	recipe := NewRecipe()

	recipe.Metadata.Name = "test"
	recipe.Metadata.Version = "v1.0.0"
	recipe.Variables = []Variable{{Name: "foo"}}
	recipe.Templates = map[string]File{
		"empty-file":                           NewFile([]byte("")),
		"empty-file-with-spaces":               NewFile([]byte(" ")),
		"empty-file-with-tabulator":            NewFile([]byte("\t")),
		"empty-file-with-spaces-and-newline-1": NewFile([]byte(" \n")),
		"empty-file-with-spaces-and-newline-2": NewFile([]byte(" \n ")),
		"file-with-empty-variable":             NewFile([]byte(" {{ .Variables.foo }} ")),
	}

	sauce, err := recipe.Execute(TestRenderEngine{}, VariableValues{"foo": ""}, uuid.Must(uuid.NewV4()))
	if err != nil {
		t.Fatalf("Failed to render recipe: %s", err)
	}

	if len(sauce.Files) > 0 {
		failingFiles := make([]string, len(sauce.Files))

		i := 0
		for k := range sauce.Files {
			failingFiles[i] = k
			i++
		}
		t.Fatalf("Rendered templates contains empty files, which exist on the output. Failing files: %s", failingFiles)
	}
}

func TestRecipeRenderWithTemplateExtension(t *testing.T) {
	recipe := NewRecipe()
	recipe.Metadata.Name = "test"
	recipe.Metadata.Version = "v1.0.0"
	recipe.Variables = []Variable{{Name: "foo"}}
	recipe.Templates = map[string]File{
		"subdirectory/file.md.tmpl": NewFile([]byte("{{ .Variables.foo }}")),
		"file":                      NewFile([]byte("{{ .Variables.foo }}")),
	}
	recipe.Metadata.TemplateExtension = ".tmpl"

	sauce, err := recipe.Execute(TestRenderEngine{}, VariableValues{"foo": "bar"}, uuid.Must(uuid.NewV4()))
	if err != nil {
		t.Fatalf("Failed to render recipe: %s", err)
	}

	if len(sauce.Files) != 2 {
		t.Fatalf("Expected 2 files, got %d", len(sauce.Files))
	}

	if string(sauce.Files["subdirectory/file.md"].Content) != "bar" {
		t.Fatalf("Expected file content to be \"bar\", got \"%s\"", sauce.Files["subdirectory/file.md"].Content)
	}

	if string(sauce.Files["file"].Content) != "{{ .Variables.foo }}" {
		t.Fatalf("Expected file content to be \"{{ .Variables.foo }}\", got \"%s\"", sauce.Files["file"].Content)
	}
}
