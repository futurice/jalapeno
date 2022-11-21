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

	if err = os.MkdirAll(filepath.Join(dir, ".jalapeno"), 0755); err != nil {
		t.Error("cannot create metadata dir", err)
	}

	recipes := `name: foo
version: v1.0.0
description: foo recipe
---
name: bar
version: v2.0.0
description: bar recipe
`

	if err = os.WriteFile(filepath.Join(dir, ".jalapeno", "recipe.yml"), []byte(recipes), 0644); err != nil {
		t.Error("cannot write recipe metadata file", err)
	}

	loaded, err := LoadRendered(dir)
	if err != nil {
		t.Error("failed to load recipes", err)
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
