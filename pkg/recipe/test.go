package recipe

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/futurice/jalapeno/pkg/engine"
	"github.com/gofrs/uuid"
	"github.com/kylelemons/godebug/diff"
)

type Test struct {
	// Name of the test case. Defined by directory name of the test case
	Name string `yaml:"-"`

	// Values to use to render the recipe templates
	Values VariableValues `yaml:"values"`

	// Expected initHelp of the recipe when rendered with the values specified in the test
	ExpectedInitHelp string `yaml:"expectedInitHelp,omitempty"`

	// If true, test will not fail if the templates generates more files than the test specifies
	IgnoreExtraFiles bool `yaml:"ignoreExtraFiles,omitempty"`

	// Snapshots of the rendered templates which were rendered with the values specified in the test
	Files map[string]File `yaml:"-"`
}

// Random hardcoded UUID
var TestID uuid.UUID = uuid.Must(uuid.FromString("12345678-1234-5678-1234-567812345678"))

var (
	ErrNoTestsSpecified    error = errors.New("no tests specified")
	ErrTestWrongFileAmount       = errors.New("recipe rendered different amount of files than expected")
	ErrTestMissingFile           = errors.New("recipe did not render file which was expected")
	ErrTestContentMismatch       = errors.New("the contents of the files did not match")
)

func (t Test) Validate() error {
	if t.Name == "" {
		return errors.New("test name can not be empty")
	}

	if err := t.Values.Validate(); err != nil {
		return fmt.Errorf("variable values were invalid: %w", err)
	}

	return nil
}

func (re *Recipe) RunTests() []error {
	errors := make([]error, len(re.Tests))
	for i, t := range re.Tests {
		sauce, err := re.Execute(engine.New(), t.Values, TestID)
		if err != nil {
			errors[i] = err
			continue
		}

		if t.ExpectedInitHelp != "" {
			actualInitHelp, err := sauce.RenderInitHelp()
			if err != nil {
				errors[i] = fmt.Errorf("could not render init help: %w", err)
			} else if actualInitHelp != t.ExpectedInitHelp {
				errors[i] = fmt.Errorf("expected init help did not match the actual init help. Expected: %s, Actual: %s", t.ExpectedInitHelp, actualInitHelp)
			}
			continue
		}

		if (t.IgnoreExtraFiles && len(t.Files) > len(sauce.Files)) || (!t.IgnoreExtraFiles && len(t.Files) != len(sauce.Files)) {
			leftOuterJoin := func(a, b map[string]File) (diff []string) {
				for key := range a {
					if _, found := b[key]; !found {
						diff = append(diff, key)
					}
				}
				return
			}

			if len(t.Files) > len(sauce.Files) {
				errors[i] = fmt.Errorf("%w: following files were missing: %s", ErrTestWrongFileAmount, leftOuterJoin(t.Files, sauce.Files))
			} else {
				errors[i] = fmt.Errorf("%w: following files were extra: %s", ErrTestWrongFileAmount, leftOuterJoin(sauce.Files, t.Files))
			}

			continue
		}

		for key, tFile := range t.Files {
			if file, ok := sauce.Files[key]; !ok {
				errors[i] = fmt.Errorf("%w: file '%s'", ErrTestMissingFile, key)
				continue
			} else {
				if !bytes.Equal(tFile.Content, file.Content) {
					errors[i] = fmt.Errorf("%w for file '%s':\n%s", ErrTestContentMismatch, key, diff.Diff(string(tFile.Content), string(file.Content)))
					continue
				}
			}
		}
	}

	return errors
}
