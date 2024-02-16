package recipeutil_test

import (
	"testing"

	"github.com/futurice/jalapeno/pkg/recipeutil"
)

func TestCreateFileTree(t *testing.T) {
	testCases := []struct {
		Name     string
		Root     string
		Input    map[string]recipeutil.FileStatus
		Expected string
	}{
		{
			"single_file_in_root",
			".",
			map[string]recipeutil.FileStatus{
				"foo": recipeutil.FileUnchanged,
			},
			`.
└── foo
`,
		},
		{
			"multiple_files_in_root",
			".",
			map[string]recipeutil.FileStatus{
				"foo": recipeutil.FileUnchanged,
				"bar": recipeutil.FileUnchanged,
			},
			`.
├── bar
└── foo
`,
		},
		{
			"nested_file",
			".",
			map[string]recipeutil.FileStatus{
				"foo/bar/baz": recipeutil.FileUnchanged,
			},
			`.
└── foo
    └── bar
        └── baz
`,
		},
		{
			"nested_files_in_same_dir",
			".",
			map[string]recipeutil.FileStatus{
				"foo/bar/baz-1": recipeutil.FileUnchanged,
				"foo/bar/baz-2": recipeutil.FileUnchanged,
			},
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
			map[string]recipeutil.FileStatus{
				"foo/bar-1/baz": recipeutil.FileUnchanged,
				"foo/bar-2/baz": recipeutil.FileUnchanged,
			},
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
			map[string]recipeutil.FileStatus{
				"b/a/b": recipeutil.FileUnchanged,
				"b/b/a": recipeutil.FileUnchanged,
				"c/a":   recipeutil.FileUnchanged,
				"a/a/c": recipeutil.FileUnchanged,
				"a/a/a": recipeutil.FileUnchanged,
				"a/a/b": recipeutil.FileUnchanged,
			},
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
		{
			"file statuses",
			".",
			map[string]recipeutil.FileStatus{
				"a": recipeutil.FileUnchanged,
				"b": recipeutil.FileAdded,
				"c": recipeutil.FileModified,
				"d": recipeutil.FileDeleted,
			},
			`.
├── a
├── b (added)
├── c (modified)
└── d (deleted)
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
