package recipe

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/futurice/jalapeno/pkg/engine"
	"github.com/gofrs/uuid"
)

type Test struct {
	// Name of the test case. Defined by directory name of the test case
	Name string `yaml:"-"`

	// Values to use to render the recipe templates
	Values VariableValues `yaml:"values"`

	// Snapshots of the rendered templates which were rendered with the values specified in the test
	Files map[string]File `yaml:"-"`

	// If true, test will not fail if the templates generates more files than the test specifies
	IgnoreExtraFiles bool `yaml:"ignoreExtraFiles"`
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

	return nil
}

func (re *Recipe) RunTests() []error {
	errors := make([]error, len(re.Tests))
	for i, t := range re.Tests {
		sauce, err := re.Execute(engine.New(), t.Values, TestID)
		if err != nil {
			errors[i] = fmt.Errorf("%w", err)
			continue
		}

		if (t.IgnoreExtraFiles && len(t.Files) > len(sauce.Files)) || (!t.IgnoreExtraFiles && len(t.Files) != len(sauce.Files)) {
			// TODO: show which files were missing/extra
			errors[i] = ErrTestWrongFileAmount
			continue
		}

		for key, tFile := range t.Files {
			if file, ok := sauce.Files[key]; !ok {
				errors[i] = fmt.Errorf("%w: file '%s'", ErrTestMissingFile, key)
				continue
			} else {
				if !bytes.Equal(tFile.Content, file.Content) {
					errors[i] = fmt.Errorf("%w: file '%s'.\nExpected:\n%s\n\nActual:\n%s", ErrTestContentMismatch, key, tFile.Content, file.Content)
					continue
				}
			}
		}
	}

	return errors
}
