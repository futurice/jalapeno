package recipe

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadNoRenderedRecipes(t *testing.T) {
	dir, err := os.MkdirTemp("", "jalapeno-test-render")
	if err != nil {
		t.Error("cannot create temp dir", err)
	}

	loaded, err := LoadRendered(dir)
	if err != nil {
		t.Error("unexpected error when loading from empty dir", err)
	}
	if len(loaded) != 0 {
		t.Error("expected slice of length 0, got", loaded)
	}
}

func TestLoadMultipleRenderedRecipes(t *testing.T) {
	dir, err := os.MkdirTemp("", "jalapeno-test-render")
	if err != nil {
		t.Error("cannot create temp dir", err)
	}

	if err = os.MkdirAll(filepath.Join(dir, RenderedRecipeDirName), 0755); err != nil {
		t.Error("cannot create metadata dir", err)
	}

	if err = os.WriteFile(filepath.Join(dir, "first.md"), []byte("# first"), 0644); err != nil {
		t.Error("cannot write rendered template", err)
	}

	if err = os.WriteFile(filepath.Join(dir, "second.md"), []byte("# second"), 0644); err != nil {
		t.Error("cannot write rendered template", err)
	}

	recipes := `apiVersion: v1
name: foo
version: v1.0.0
description: foo recipe
files:
  first.md:
    checksum: sha256:a04042ce4a5e66443c5a26ef2d4432aa535421286c062ea7bf55cba5bae15ef4
---
apiVersion: v1
name: bar
version: v2.0.0
description: bar recipe
files:
  second.md:
    checksum: sha256:1b42293a96dbdcf36ee77dcbee6e2e2804ab085d32e6a2de7736198a0d111044
`

	if err = os.WriteFile(filepath.Join(dir, RenderedRecipeDirName, RecipeFileName+YAMLExtension), []byte(recipes), 0644); err != nil {
		t.Error("cannot write recipe metadata file", err)
	}

	loaded, err := LoadRendered(dir)
	if err != nil {
		t.Errorf("failed to load recipes: %s", err)
	}
	if len(loaded) != 2 {
		t.Errorf("expected to load 2 recipes, loaded %d", len(loaded))
	}

	if loaded[0].Name != "foo" {
		t.Errorf("expected 'foo' as the first recipe name, got %s", loaded[0].Name)
	}
	if loaded[1].Name != "bar" {
		t.Errorf("expected 'bar' as the first recipe name, got %s", loaded[0].Name)
	}
}
