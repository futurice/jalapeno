package recipeutil_test

import (
	"testing"

	"github.com/futurice/jalapeno/pkg/recipeutil"
)

func TestCreateFileTree(t *testing.T) {
	testCases := []struct {
		Name     string
		Root     string
		Input    []string
		Expected string
	}{
		{
			"single_file_in_root",
			".",
			[]string{"foo"},
			`.
└── foo
`,
		},
		{
			"multiple_files_in_root",
			".",
			[]string{"foo", "bar"},
			`.
├── bar
└── foo
`,
		},
		{
			"nested_file",
			".",
			[]string{"foo/bar/baz"},
			`.
└── foo
    └── bar
        └── baz
`,
		},
		{
			"nested_files_in_same_dir",
			".",
			[]string{"foo/bar/baz-1", "foo/bar/baz-2"},
			`.
└── foo
    └── bar
        ├── baz-1
        └── baz-2
`,
		},
		{
			"nested_files_in_different_dirs",
			".",
			[]string{"foo/bar-1/baz", "foo/bar-2/baz"},
			`.
└── foo
    ├── bar-1
    │   └── baz
    └── bar-2
        └── baz
`,
		},
		{
			"files_are_alphabetically_sorted",
			".",
			[]string{"b/a/b", "b/b/a", "c/a", "a/a/c", "a/a/a", "a/a/b"},
			`.
├── a
│   └── a
│       ├── a
│       ├── b
│       └── c
├── b
│   ├── a
│   │   └── b
│   └── b
│       └── a
└── c
    └── a
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			actual := recipeutil.CreateFileTree(tc.Root, tc.Input)
			if actual != tc.Expected {
				t.Errorf("expected:\n%s\nactual:\n%s", tc.Expected, actual)
			}
		})
	}
}
