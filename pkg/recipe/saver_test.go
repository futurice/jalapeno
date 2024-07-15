package recipe

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/futurice/jalapeno/pkg/engine"
	"github.com/gofrs/uuid"
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
	re.Templates = map[string]File{
		"foo.md":         NewFile([]byte("foo")),
		"foo/bar.md":     NewFile([]byte("bar")),
		"foo/bar/baz.md": NewFile([]byte("baz")),
	}
	re.Tests = []Test{
		{
			Name: "baz_test",
			Files: map[string]File{
				"foo.md":         NewFile([]byte("foo")),
				"foo/bar.md":     NewFile([]byte("bar")),
				"foo/bar/baz.md": NewFile([]byte("baz")),
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
		filepath.Join(RecipeFileName + YAMLExtension),
		filepath.Join("templates", "foo.md"),
		filepath.Join("templates", "foo", "bar.md"),
		filepath.Join("templates", "foo", "bar", "baz.md"),
		filepath.Join("tests", re.Tests[0].Name, RecipeTestMetaFileName+YAMLExtension),
		filepath.Join("tests", re.Tests[0].Name, RecipeTestFilesDirName, "foo.md"),
		filepath.Join("tests", re.Tests[0].Name, RecipeTestFilesDirName, "foo", "bar.md"),
		filepath.Join("tests", re.Tests[0].Name, RecipeTestFilesDirName, "foo", "bar", "baz.md"),
	}

	checkFiles(t, dir, expectedFiles)
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
	re.Templates = map[string]File{
		"foo.md":         NewFile([]byte("foo")),
		"foo/bar.md":     NewFile([]byte("bar")),
		"foo/bar/baz.md": NewFile([]byte("baz")),
	}
	re.Tests = []Test{
		{
			Name: "baz_test",
			Files: map[string]File{
				"foo/bar/baz.md": NewFile([]byte("baz")),
			},
		},
	}

	err = re.Validate()
	if err != nil {
		t.Fatalf("test recipe was not valid: %s", err)
	}

	sauce, err := re.Execute(engine.New(), nil, uuid.Must(uuid.NewV4()))
	if err != nil {
		t.Fatalf("recipe execution failed: %s", err)
	}

	err = sauce.Save(dir)
	if err != nil {
		t.Fatalf("failed to save sauce: %s", err)
	}

	expectedFiles := []string{
		filepath.Join(SauceDirName, SaucesFileName+YAMLExtension),
		filepath.Join("foo.md"),
		filepath.Join("foo", "bar.md"),
		filepath.Join("foo", "bar", "baz.md"),
	}

	checkFiles(t, dir, expectedFiles)
}

func TestSaveSauceDoesNotWriteOutsideDest(t *testing.T) {
	dir, err := os.MkdirTemp("", "jalapeno-test-saver")
	if err != nil {
		t.Fatalf("cannot create temp dir: %s", err)
	}
	defer os.RemoveAll(dir)

	re := NewRecipe()
	re.Name = "Test"
	re.Version = "v0.0.1"
	re.Templates = map[string]File{
		"../foo.md": NewFile([]byte("foo")),
	}

	err = re.Validate()
	if err != nil {
		t.Fatalf("test recipe was not valid: %s", err)
	}

	sauce, err := re.Execute(engine.New(), nil, uuid.Must(uuid.NewV4()))
	if err != nil {
		t.Fatalf("recipe execution failed: %s", err)
	}

	err = sauce.Save(dir)
	if err == nil {
		t.Fatalf("should not have saved sauce")
	}

	if !strings.Contains(err.Error(), "file path escapes destination") {
		t.Fatalf("error received was not expected: %s", err)
	}
}

func checkFiles(t *testing.T, dir string, expectedFiles []string) {
	files := make(map[string]bool, len(expectedFiles))
	for _, file := range expectedFiles {
		files[filepath.Join(dir, file)] = false
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if _, found := files[path]; !found {
			t.Fatalf("unexpected file '%s' found", path)
		} else {
			files[path] = true
		}

		return nil
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	for file, found := range files {
		if !found {
			t.Fatalf("expected file '%s' not found", file)
		}
	}
}
