package recipe

import (
	"errors"
	"testing"
)

type TestWithExpectedOutcome struct {
	Test
	ExpectedError error
}

func TestRecipeTests(t *testing.T) {
	scenarios := []struct {
		name      string
		templates map[string]File
		tests     []TestWithExpectedOutcome
	}{
		{
			"single_pass",
			map[string]File{
				"foo.txt":                NewFile([]byte("foo")),
				"var.txt":                NewFile([]byte("{{ .Variables.VAR }}")),
				"var_with_func_pipe.txt": NewFile([]byte("{{ sha1sum .Variables.VAR }}")),
			},
			[]TestWithExpectedOutcome{
				{
					Test: Test{
						Values: VariableValues{
							"VAR": "var",
						},
						Files: map[string]File{
							"foo.txt":                NewFile([]byte("foo")),
							"var.txt":                NewFile([]byte("var")),
							"var_with_func_pipe.txt": NewFile([]byte("e5b4e786e382d03c28e9edfab2d8149378ae69df")), // echo -n "var" | shasum -a 1
						},
					},
					ExpectedError: nil,
				},
			},
		},
		{
			"one_pass_one_failing",
			map[string]File{
				"foo.txt": NewFile([]byte("foo")),
			},
			[]TestWithExpectedOutcome{
				{
					Test: Test{
						Files: map[string]File{
							"foo.txt": NewFile([]byte("foo")),
						},
					},
					ExpectedError: nil,
				},
				{
					Test: Test{
						Files: map[string]File{
							"foo.txt": NewFile([]byte("bar")),
						},
					},
					ExpectedError: ErrTestContentMismatch,
				},
			},
		},
		{
			"no_tests",
			map[string]File{},
			nil,
		},
		{
			"content_mismatch",
			map[string]File{
				"foo.txt": NewFile([]byte("bar")),
			},
			[]TestWithExpectedOutcome{
				{
					Test: Test{
						Files: map[string]File{
							"foo.txt": NewFile([]byte("foo")),
						},
					},
					ExpectedError: ErrTestContentMismatch,
				},
			},
		},
		{
			"expected_more_files",
			map[string]File{
				"foo.txt": NewFile([]byte("foo")),
			},
			[]TestWithExpectedOutcome{
				{
					Test: Test{
						Files: map[string]File{
							"foo.txt": NewFile([]byte("foo")),
							"bar.txt": NewFile([]byte("bar")),
						},
					},
					ExpectedError: ErrTestWrongFileAmount,
				},
			},
		},
		{
			"expected_less_files",
			map[string]File{
				"foo.txt": NewFile([]byte("foo")),
				"bar.txt": NewFile([]byte("bar")),
			},
			[]TestWithExpectedOutcome{
				{
					Test: Test{
						Files: map[string]File{
							"foo.txt": NewFile([]byte("foo")),
						},
					},
					ExpectedError: ErrTestWrongFileAmount,
				},
			},
		},
		{
			"unexpected_file_rendered",
			map[string]File{
				"foo.txt": NewFile([]byte("foo")),
				"baz.txt": NewFile([]byte("baz")),
			},
			[]TestWithExpectedOutcome{
				{
					Test: Test{
						Files: map[string]File{
							"foo.txt": NewFile([]byte("foo")),
							"bar.txt": NewFile([]byte("bar")),
						},
					},
					ExpectedError: ErrTestMissingFile,
				},
			},
		},
		{
			"skip_extra_files",
			map[string]File{
				"foo.txt":   NewFile([]byte("foo")),
				"bar.txt":   NewFile([]byte("bar")),
				"extra.txt": NewFile([]byte("extra")),
			},
			[]TestWithExpectedOutcome{
				{
					Test: Test{
						Files: map[string]File{
							"foo.txt": NewFile([]byte("foo")),
							"bar.txt": NewFile([]byte("bar")),
						},
						IgnoreExtraFiles: true,
					},
					ExpectedError: nil,
				},
			},
		},
		{
			"missing_file_with_skip_extra_files_enabled",
			map[string]File{
				"foo.txt": NewFile([]byte("foo")),
				"bar.txt": NewFile([]byte("bar")),
			},
			[]TestWithExpectedOutcome{
				{
					Test: Test{
						Files: map[string]File{
							"foo.txt":   NewFile([]byte("foo")),
							"bar.txt":   NewFile([]byte("bar")),
							"extra.txt": NewFile([]byte("extra")),
						},
						IgnoreExtraFiles: true,
					},
					ExpectedError: ErrTestWrongFileAmount,
				},
			},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(tt *testing.T) {
			tests := make([]Test, len(scenario.tests))
			for i, t := range scenario.tests {
				tests[i] = t.Test
				tests[i].Name = scenario.name
			}
			recipe := NewRecipe()
			recipe.Name = "foo"
			recipe.Version = "v0.0.0"
			recipe.Tests = tests
			recipe.Templates = scenario.templates

			errs := recipe.RunTests()
			for i, err := range errs {
				expectedErr := scenario.tests[i].ExpectedError
				if expectedErr == nil && err != nil {
					tt.Fatalf("Tests returned error when not expected: %s", err)
				} else if expectedErr != nil && err == nil {
					tt.Fatalf("Tests did not return error when expected. Expected error: %s", expectedErr)
				}

				if err != nil && !errors.Is(err, expectedErr) {
					tt.Fatalf("Tests did not return expected error. Expected: '%s'. Actual: '%s'", expectedErr, err)
				}
			}
		})
	}
}
