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
				"foo": recipeutil.FileUnknown,
			},
			`.
└── foo
`,
		},
		{
			"multiple_files_in_root",
			".",
			map[string]recipeutil.FileStatus{
				"foo": recipeutil.FileUnknown,
				"bar": recipeutil.FileUnknown,
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
				"foo/bar/baz": recipeutil.FileUnknown,
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
				"foo/bar/baz-1": recipeutil.FileUnknown,
				"foo/bar/baz-2": recipeutil.FileUnknown,
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
				"foo/bar-1/baz": recipeutil.FileUnknown,
				"foo/bar-2/baz": recipeutil.FileUnknown,
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
				"b/a/b": recipeutil.FileUnknown,
				"b/b/a": recipeutil.FileUnknown,
				"c/a":   recipeutil.FileUnknown,
				"a/a/c": recipeutil.FileUnknown,
				"a/a/a": recipeutil.FileUnknown,
				"a/a/b": recipeutil.FileUnknown,
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
			"directory_as_entry",
			".",
			map[string]recipeutil.FileStatus{
				"a/": recipeutil.FileUnknown,
			},
			`.
└── a/
`,
		},
		{
			"file statuses",
			".",
			map[string]recipeutil.FileStatus{
				"a": recipeutil.FileUnknown,
				"b": recipeutil.FileUnchanged,
				"c": recipeutil.FileAdded,
				"d": recipeutil.FileModified,
				"e": recipeutil.FileDeleted,
			},
			`.
├── a
├── b (unchanged)
├── c (added)
├── d (modified)
└── e (deleted)
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
