package recipe

import (
	"errors"
	"testing"
)

func TestRecipeTests(t *testing.T) {
	tests := []struct {
		name      string
		templates map[string][]byte
		tests     []Test
		expectErr error
	}{
		{
			"pass",
			map[string][]byte{
				"foo.txt": []byte("foo"),
				"var.txt": []byte("{{ .Variables.VAR }}"),
			},
			[]Test{
				{
					Values: VariableValues{
						"VAR": "var",
					},
					Files: map[string][]byte{
						"foo.txt": []byte("foo"),
						"var.txt": []byte("var"),
					},
				},
			},
			nil,
		},
		{
			"no_tests",
			map[string][]byte{},
			nil,
			ErrNoTestsSpecified,
		},
		{
			"content_mismatch",
			map[string][]byte{
				"foo.txt": []byte("bar"),
			},
			[]Test{
				{
					Files: map[string][]byte{
						"foo.txt": []byte("foo"),
					},
				},
			},
			ErrTestContentMismatch,
		},
		{
			"expected_more_files",
			map[string][]byte{
				"foo.txt": []byte("foo"),
			},
			[]Test{
				{
					Files: map[string][]byte{
						"foo.txt": []byte("foo"),
						"bar.txt": []byte("bar"),
					},
				},
			},
			ErrTestWrongFileAmount,
		},
		{
			"expected_less_files",
			map[string][]byte{
				"foo.txt": []byte("foo"),
				"bar.txt": []byte("bar"),
			},
			[]Test{
				{
					Files: map[string][]byte{
						"foo.txt": []byte("foo"),
					},
				},
			},
			ErrTestWrongFileAmount,
		},
		{
			"unexpected_file_rendered",
			map[string][]byte{
				"foo.txt": []byte("foo"),
				"baz.txt": []byte("baz"),
			},
			[]Test{
				{
					Files: map[string][]byte{
						"foo.txt": []byte("foo"),
						"bar.txt": []byte("bar"),
					},
				},
			},
			ErrTestMissingFile,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recipe := new()
			recipe.tests = tt.tests
			recipe.Templates = tt.templates

			err := recipe.RunTests()
			if tt.expectErr == nil && err != nil {
				t.Errorf("Tests returned error when not expected: %s", err)
				return
			} else if tt.expectErr != nil && err == nil {
				t.Errorf("Tests did not return error when expected. Expected error: %s", tt.expectErr)
				return
			}

			if err != nil && !errors.Is(err, tt.expectErr) {
				t.Errorf("Tests did not return expected error. Expected: '%s'. Actual: '%s'", tt.expectErr, err)
				return
			}
		})
	}
}
