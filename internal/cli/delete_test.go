package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/gofrs/uuid"
)

func TestDelete(t *testing.T) {
	// Create a temporary directory for testing
	dir, err := os.MkdirTemp("", "jalapeno-test-delete")
	if err != nil {
		t.Fatalf("cannot create temp dir: %s", err)
	}
	defer os.RemoveAll(dir)

	// Create test files and directories
	if err = os.MkdirAll(filepath.Join(dir, recipe.SauceDirName), 0755); err != nil {
		t.Fatalf("cannot create metadata dir: %s", err)
	}

	// Create some test files that will be "managed" by the sauces
	testFiles := []string{"first.md", "second.md"}
	for _, f := range testFiles {
		if err = os.WriteFile(filepath.Join(dir, f), []byte("# "+f), 0644); err != nil {
			t.Fatalf("cannot write test file: %s", err)
		}
	}

	// Create test sauces
	id1 := uuid.Must(uuid.NewV4())
	id2 := uuid.Must(uuid.NewV4())

	sauces := []*recipe.Sauce{
		{
			APIVersion: "v1",
			ID:         id1,
			Recipe: recipe.Recipe{
				Metadata: recipe.Metadata{
					APIVersion: "v1",
					Name:      "foo",
					Version:   "v1.0.0",
				},
			},
			Files: map[string]recipe.File{
				"first.md": recipe.NewFile([]byte("# first")),
			},
		},
		{
			APIVersion: "v1",
			ID:         id2,
			Recipe: recipe.Recipe{
				Metadata: recipe.Metadata{
					APIVersion: "v1",
					Name:      "bar",
					Version:   "v2.0.0",
				},
			},
			Files: map[string]recipe.File{
				"second.md": recipe.NewFile([]byte("# second")),
			},
		},
	}

	if err = recipe.SaveSauces(dir, sauces); err != nil {
		t.Fatalf("cannot save test sauces: %s", err)
	}

	t.Run("delete specific sauce", func(t *testing.T) {
		cmd := NewRootCmd()
		cmd.SetArgs([]string{"delete", id1.String(), "--dir", dir})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("failed to execute delete command: %s", err)
		}

		// Check that first.md was deleted
		if _, err := os.Stat(filepath.Join(dir, "first.md")); !os.IsNotExist(err) {
			t.Error("first.md should have been deleted")
		}

		// Check that second.md still exists
		if _, err := os.Stat(filepath.Join(dir, "second.md")); err != nil {
			t.Error("second.md should still exist")
		}

		// Check that only one sauce remains
		remainingSauces, err := recipe.LoadSauces(dir)
		if err != nil {
			t.Fatalf("failed to load sauces: %s", err)
		}

		if len(remainingSauces) != 1 {
			t.Errorf("expected 1 sauce, got %d", len(remainingSauces))
		}

		if remainingSauces[0].ID != id2 {
			t.Error("wrong sauce was deleted")
		}
	})

	t.Run("delete all sauces", func(t *testing.T) {
		cmd := NewRootCmd()
		cmd.SetArgs([]string{"delete", "--all", "--dir", dir})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("failed to execute delete command: %s", err)
		}

		// Check that both files were deleted
		for _, f := range testFiles {
			if _, err := os.Stat(filepath.Join(dir, f)); !os.IsNotExist(err) {
				t.Errorf("%s should have been deleted", f)
			}
		}

		// Check that .jalapeno directory was deleted
		if _, err := os.Stat(filepath.Join(dir, recipe.SauceDirName)); !os.IsNotExist(err) {
			t.Error(".jalapeno directory should have been deleted")
		}
	})

	t.Run("delete with invalid sauce ID", func(t *testing.T) {
		cmd := NewRootCmd()
		cmd.SetArgs([]string{"delete", "invalid-uuid", "--dir", dir})

		if err := cmd.Execute(); err == nil {
			t.Fatal("expected error with invalid sauce ID")
		}
	})

	t.Run("delete without sauce ID or --all flag", func(t *testing.T) {
		cmd := NewRootCmd()
		cmd.SetArgs([]string{"delete", "--dir", dir})

		if err := cmd.Execute(); err == nil {
			t.Fatal("expected error when no sauce ID or --all flag provided")
		}
	})

	t.Run("delete non-existent sauce", func(t *testing.T) {
		nonExistentID := uuid.Must(uuid.NewV4())
		cmd := NewRootCmd()
		cmd.SetArgs([]string{"delete", nonExistentID.String(), "--dir", dir})

		if err := cmd.Execute(); err == nil {
			t.Fatal("expected error when deleting non-existent sauce")
		}
	})
} 