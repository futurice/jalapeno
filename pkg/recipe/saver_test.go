package recipe

import (
	"os"
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

	// TODO: Check output files
}
