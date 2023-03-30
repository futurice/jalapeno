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
				"foo.txt":                []byte("foo"),
				"var.txt":                []byte("{{ .Variables.VAR }}"),
				"var_with_func_pipe.txt": []byte("{{ sha1sum .Variables.VAR }}"),
			},
			[]Test{
				{
					Values: VariableValues{
						"VAR": "var",
					},
					Files: map[string]TestFile{
						"foo.txt":                []byte("foo"),
						"var.txt":                []byte("var"),
						"var_with_func_pipe.txt": []byte("e5b4e786e382d03c28e9edfab2d8149378ae69df"), // echo -n "var" | shasum -a 1
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
					Files: map[string]TestFile{
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
					Files: map[string]TestFile{
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
					Files: map[string]TestFile{
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
					Files: map[string]TestFile{
						"foo.txt": []byte("foo"),
						"bar.txt": []byte("bar"),
					},
				},
			},
			ErrTestMissingFile,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			recipe := new()
			recipe.tests = test.tests
			recipe.Templates = test.templates

			err := recipe.RunTests()
			if test.expectErr == nil && err != nil {
				tt.Fatalf("Tests returned error when not expected: %s", err)
			} else if test.expectErr != nil && err == nil {
				tt.Fatalf("Tests did not return error when expected. Expected error: %s", test.expectErr)
			}

			if err != nil && !errors.Is(err, test.expectErr) {
				tt.Fatalf("Tests did not return expected error. Expected: '%s'. Actual: '%s'", test.expectErr, err)
			}
		})
	}
}
