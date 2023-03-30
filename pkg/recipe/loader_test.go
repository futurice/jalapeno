package recipe

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadNoRenderedRecipes(t *testing.T) {
	dir, err := os.MkdirTemp("", "jalapeno-test-loader")
	if err != nil {
		t.Fatalf("cannot create temp dir: %s", err)
	}
	defer os.RemoveAll(dir)

	loaded, err := LoadRendered(dir)
	if err != nil {
		t.Fatalf("unexpected error when loading from empty dir: %s", err)
	}
	if len(loaded) != 0 {
		t.Fatalf("expected slice of length 0, got %d", len(loaded))
	}
}

func TestLoadMultipleRenderedRecipes(t *testing.T) {
	dir, err := os.MkdirTemp("", "jalapeno-test-loader")
	if err != nil {
		t.Fatalf("cannot create temp dir: %s", err)
	}
	defer os.RemoveAll(dir)

	if err = os.MkdirAll(filepath.Join(dir, RenderedRecipeDirName), 0755); err != nil {
		t.Fatalf("cannot create metadata dir: %s", err)
	}

	if err = os.WriteFile(filepath.Join(dir, "first.md"), []byte("# first"), 0644); err != nil {
		t.Fatalf("cannot write rendered template: %s", err)
	}

	if err = os.WriteFile(filepath.Join(dir, "second.md"), []byte("# second"), 0644); err != nil {
		t.Fatalf("cannot write rendered template: %s", err)
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
		t.Fatalf("cannot write recipe metadata file: %s", err)
	}

	loaded, err := LoadRendered(dir)
	if err != nil {
		t.Fatalf("failed to load recipes: %s", err)
	}
	if len(loaded) != 2 {
		t.Fatalf("expected to load 2 recipes, loaded %d", len(loaded))
	}

	if loaded[0].Name != "foo" {
		t.Fatalf("expected 'foo' as the first recipe name, got %s", loaded[0].Name)
	}
	if loaded[1].Name != "bar" {
		t.Fatalf("expected 'bar' as the first recipe name, got %s", loaded[0].Name)
	}
}

func TestLoadTests(t *testing.T) {
	dir, err := os.MkdirTemp("", "jalapeno-test-loader")
	if err != nil {
		t.Fatalf("cannot create temp dir: %s", err)
	}
	defer os.RemoveAll(dir)

	if err = os.MkdirAll(filepath.Join(dir, RecipeTemplatesDirName), 0755); err != nil {
		t.Fatalf("cannot create templates dir: %s", err)
	}

	contents := "# file"
	if err = os.WriteFile(filepath.Join(dir, RecipeTemplatesDirName, "file.md"), []byte(contents), 0644); err != nil {
		t.Fatalf("cannot write rendered template: %s", err)
	}

	recipe := `apiVersion: v1
name: foo
version: v1.0.0
description: foo recipe
`

	if err = os.WriteFile(filepath.Join(dir, RecipeFileName+YAMLExtension), []byte(recipe), 0644); err != nil {
		t.Fatal("cannot write recipe file", err)
	}

	testCase := `values: {}
files:
  "file.md": IyBmaWxl
`
	if err = os.MkdirAll(filepath.Join(dir, RecipeTestsDirName), 0755); err != nil {
		t.Fatalf("cannot create test dir: %s", err)
	}

	if err = os.WriteFile(filepath.Join(dir, RecipeTestsDirName, "test_foo"+YAMLExtension), []byte(testCase), 0644); err != nil {
		t.Fatalf("cannot write recipe test file: %s", err)
	}

	loaded, err := Load(dir)
	if err != nil {
		t.Fatalf("failed to load the recipe: %s", err)
	}

	if len(loaded.tests) != 1 {
		t.Fatal("failed to load recipe tests")
	}

	if !bytes.Equal(loaded.tests[0].Files["file.md"], []byte(contents)) {
		t.Fatalf("loader did not decode recipe test files correctly, expected %s, actual %s", contents, loaded.tests[0].Files["file.md"])
	}
}
