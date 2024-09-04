package diff_test

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/futurice/jalapeno/pkg/diff"
)

func TestGetUnifiedDiffLines(t *testing.T) {
	testCases := []struct {
		name                string
		a                   string
		b                   string
		expectedIsDifferent bool
		expectedLines       []string
	}{
		{
			name:                "both_empty",
			a:                   "",
			b:                   "",
			expectedIsDifferent: false,
			expectedLines:       nil,
		},
		{
			name:                "both_empty_terminated",
			a:                   "\n",
			b:                   "\n",
			expectedIsDifferent: false,
			expectedLines:       nil,
		},
		{
			name:                "equal_unterminated",
			a:                   "1\n2\n3",
			b:                   "1\n2\n3",
			expectedIsDifferent: false,
			expectedLines:       nil,
		},
		{
			name:                "equal_terminated",
			a:                   "1\n2\n3\n",
			b:                   "1\n2\n3\n",
			expectedIsDifferent: false,
			expectedLines:       nil,
		},
		{
			name:                "line_added_unterminated",
			a:                   "",
			b:                   "b",
			expectedIsDifferent: true,
			// FIXME []string{"+b"} would make more sense?
			expectedLines: []string{"-", "+b"},
		},
		{
			name:                "line_added_terminated",
			a:                   "",
			b:                   "b\n",
			expectedIsDifferent: true,
			expectedLines:       []string{"+b", " "},
		},
		{
			name:                "text_added_to_line_terminated",
			a:                   "\n",
			b:                   "b\n",
			expectedIsDifferent: true,
			// FIXME totally bogus result? should be []string{"-", "+b"}
			expectedLines: []string{"+b", " ", "-"},
		},
		{
			name:                "line_changed_unterminated",
			a:                   "a",
			b:                   "b",
			expectedIsDifferent: true,
			expectedLines:       []string{"-a", "+b"},
		},
		{
			name:                "line_changed_terminated",
			a:                   "a\n",
			b:                   "b\n",
			expectedIsDifferent: true,
			expectedLines:       []string{"-a", "+b", " "},
		},
		{
			name:                "line_removed_unterminated",
			a:                   "a",
			b:                   "",
			expectedIsDifferent: true,
			// FIXME []string{"-a"} would make more sense?
			expectedLines: []string{"-a", "+"},
		},
		{
			name:                "line_removed_terminated",
			a:                   "a\n",
			b:                   "",
			expectedIsDifferent: true,
			// FIXME []string{"-a"} would make more sense?
			expectedLines: []string{"-a", " "},
		},
		{
			name:                "one_line_changed_to_two",
			a:                   "a\n",
			b:                   "b\nb\n",
			expectedIsDifferent: true,
			expectedLines:       []string{"-a", "+b", "+b", " "},
		},
		{
			name:                "two_lines_changed_to_one",
			a:                   "a\na\n",
			b:                   "b\n",
			expectedIsDifferent: true,
			expectedLines:       []string{"-a", "-a", "+b", " "},
		},
		{
			name:                "two_lines_changed",
			a:                   "a\na\n",
			b:                   "b\nb\n",
			expectedIsDifferent: true,
			expectedLines:       []string{"-a", "-a", "+b", "+b", " "},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := diff.New(tc.a, tc.b)
			lines := d.GetUnifiedDiffLines()

			if d.IsDifferent() != tc.expectedIsDifferent {
				if tc.expectedIsDifferent {
					t.Errorf("Expected texts to be different, but were the same")
				} else {
					t.Errorf("Expected texts to be the same, but were different")
				}
			}

			if !slices.Equal(lines, tc.expectedLines) {
				t.Errorf("Expected lines: [%v]; got lines: [%v]", strings.Join(tc.expectedLines, ","), strings.Join(lines, ","))
			}
		})
	}
}

func TestGetUnifiedDiff(t *testing.T) {
	// not testing all the cases as we know this is just a simple concatenation of lines produced by GetUnifiedDiffLines
	testCases := []struct {
		name                string
		a                   string
		b                   string
		expectedIsDifferent bool
		expectedDiff        string
	}{
		{
			name:                "equal",
			a:                   "a\na\n",
			b:                   "a\na\n",
			expectedIsDifferent: false,
			expectedDiff:        "",
		},
		{
			name:                "unequal",
			a:                   "a\na\n",
			b:                   "b\nb\n",
			expectedIsDifferent: true,
			expectedDiff:        "-a\n-a\n+b\n+b\n \n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := diff.New(tc.a, tc.b)

			if d.IsDifferent() != tc.expectedIsDifferent {
				if tc.expectedIsDifferent {
					t.Errorf("Expected texts to be different, but were the same")
				} else {
					t.Errorf("Expected texts to be the same, but were different")
				}
			}

			unifiedDiff := d.GetUnifiedDiff()

			if unifiedDiff != tc.expectedDiff {
				t.Errorf("Expected unified diff [%s]; got [%s]", tc.expectedDiff, unifiedDiff)
			}
		})
	}
}

func TestGetConflictResolutionTemplate(t *testing.T) {
	// not testing all the cases as we know this is just a simple concatenation
	testCases := []struct {
		name                string
		a                   string
		b                   string
		expectedIsDifferent bool
		expectedTemplate    string
	}{
		{
			name:                "both_empty",
			a:                   "",
			b:                   "",
			expectedIsDifferent: false,
			expectedTemplate:    "",
		},
		{
			name:                "equal",
			a:                   "a\na\n",
			b:                   "a\na\n",
			expectedIsDifferent: false,
			expectedTemplate:    "",
		},
		{
			name:                "unequal",
			a:                   "a\na\n",
			b:                   "b\na\n",
			expectedIsDifferent: true,
			// this is not the shortest possible obviously, but it's what the underlying library can give us
			expectedTemplate: `<<<<<<< Deleted
=======
b
>>>>>>> Added
a
<<<<<<< Deleted
a
=======
>>>>>>> Added
`,
		},
		{
			name:                "unequal_with_leading_text",
			a:                   "c\na\na\n",
			b:                   "c\nb\na\n",
			expectedIsDifferent: true,
			expectedTemplate: `c
<<<<<<< Deleted
=======
b
>>>>>>> Added
a
<<<<<<< Deleted
a
=======
>>>>>>> Added
`,
		},
		{
			name:                "unequal_with_trailing_text",
			a:                   "a\na\nc\n",
			b:                   "b\na\nc\n",
			expectedIsDifferent: true,
			expectedTemplate: `<<<<<<< Deleted
=======
b
>>>>>>> Added
a
<<<<<<< Deleted
a
=======
>>>>>>> Added
c
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := diff.New(tc.a, tc.b)

			if d.IsDifferent() != tc.expectedIsDifferent {
				if tc.expectedIsDifferent {
					t.Errorf("Expected texts to be different, but were the same")
				} else {
					t.Errorf("Expected texts to be the same, but were different")
				}
			}

			resolutionTemplate := d.GetConflictResolutionTemplate()

			if resolutionTemplate != tc.expectedTemplate {
				t.Errorf("Expected conflict resolution template [%s]; got [%s]", tc.expectedTemplate, resolutionTemplate)
			}
		})
	}
}

func TestLargerTestCases(t *testing.T) {
	testCases := []struct {
		name string
	}{
		{
			name: "case1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stringA, err := readStringFromFile(fmt.Sprintf("testdata/%s_a.txt", tc.name))
			if err != nil {
				t.Fatalf("Failed to read file a: %v", err)
			}

			stringB, err := readStringFromFile(fmt.Sprintf("testdata/%s_b.txt", tc.name))
			if err != nil {
				t.Fatalf("Failed to read file b: %v", err)
			}

			expectedDiff, err := readStringFromFile(fmt.Sprintf("testdata/%s_diff.txt", tc.name))
			if err != nil {
				t.Fatalf("Failed to read expected diff: %v", err)
			}

			expectedTemplate, err := readStringFromFile(fmt.Sprintf("testdata/%s_template.txt", tc.name))
			if err != nil {
				t.Fatalf("Failed to read expected template: %v", err)
			}

			d := diff.New(stringA, stringB)

			// warning: the expected diff files has lines which consist of a single space; some editors / IDEs
			// might be configured to strip trailing whitespace from lines, which might corrupt the file.
			unifiedDiff := d.GetUnifiedDiff()
			if unifiedDiff != expectedDiff {
				t.Errorf("Expected unified diff [%s]; got [%s]", expectedDiff, unifiedDiff)
			}

			template := d.GetConflictResolutionTemplate()
			if template != expectedTemplate {
				t.Errorf("Expected conflict resolution template [%s]; got [%s]", expectedTemplate, template)
			}
		})
	}
}

func readStringFromFile(name string) (string, error) {
	fileBytes1, err := os.ReadFile(name)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return string(fileBytes1), nil
}
