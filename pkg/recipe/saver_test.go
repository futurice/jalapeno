package recipe

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveRecipe(t *testing.T) {
	dir, err := os.MkdirTemp("", "jalapeno-test-saver")
	if err != nil {
		t.Fatalf("cannot create temp dir: %s", err)
	}
	defer os.RemoveAll(dir)

	re := NewRecipe()
	re.Name = "Test"
	re.Version = "v0.0.1"
	re.Templates = map[string][]byte{
		"foo.md":         []byte("foo"),
		"foo/bar.md":     []byte("bar"),
		"foo/bar/baz.md": []byte("baz"),
	}
	re.Tests = []Test{
		{
			Name: "baz_test",
			Files: map[string]TestFile{
				"foo/bar/baz.md": []byte("baz"),
			},
		},
	}

	err = re.Validate()
	if err != nil {
		t.Fatalf("test recipe was not valid: %s", err)
	}

	err = re.Save(dir)
	if err != nil {
		t.Fatalf("failed to save recipe: %s", err)
	}

	expectedFiles := []string{
		filepath.Join(dir, re.Name, RecipeFileName+YAMLExtension),
		filepath.Join(dir, re.Name, "templates", "foo.md"),
		filepath.Join(dir, re.Name, "templates", "foo", "bar.md"),
		filepath.Join(dir, re.Name, "templates", "foo", "bar", "baz.md"),
		filepath.Join(dir, re.Name, "tests", re.Tests[0].Name+YAMLExtension),
	}

	// TODO: check that these are _only_ files existing
	for _, expectedFile := range expectedFiles {
		if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
			t.Fatalf("expected file '%s' did not exist", expectedFile)
		}
	}
}

func TestSaveSauce(t *testing.T) {
	dir, err := os.MkdirTemp("", "jalapeno-test-saver")
	if err != nil {
		t.Fatalf("cannot create temp dir: %s", err)
	}
	defer os.RemoveAll(dir)

	re := NewRecipe()
	re.Name = "Test"
	re.Version = "v0.0.1"
	re.Templates = map[string][]byte{
		"foo.md":         []byte("foo"),
		"foo/bar.md":     []byte("bar"),
		"foo/bar/baz.md": []byte("baz"),
	}
	re.Tests = []Test{
		{
			Name: "baz_test",
			Files: map[string]TestFile{
				"foo/bar/baz.md": []byte("baz"),
			},
		},
	}

	err = re.Validate()
	if err != nil {
		t.Fatalf("test recipe was not valid: %s", err)
	}

	sauce, err := re.Execute(nil, ExecuteOptions{})
	if err != nil {
		t.Fatalf("recipe execution failed: %s", err)
	}

	err = sauce.Save(dir)
	if err != nil {
		t.Fatalf("failed to save sauce: %s", err)
	}

	expectedFiles := []string{
		filepath.Join(dir, SauceDirName, SaucesFileName+YAMLExtension),
		filepath.Join(dir, "foo.md"),
		filepath.Join(dir, "foo", "bar.md"),
		filepath.Join(dir, "foo", "bar", "baz.md"),
	}

	// TODO: check that these are _only_ files existing
	for _, expectedFile := range expectedFiles {
		if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
			t.Fatalf("expected file '%s' did not exist", expectedFile)
		}
	}
}
