package recipe

import (
	"bytes"
	"errors"
	"fmt"
)

type Test struct {
	Name   string            `yaml:"name,omitempty"`
	Values VariableValues    `yaml:"values"`
	Files  map[string][]byte `yaml:"files"`
}

var (
	ErrNoTestsSpecified    error = errors.New("no tests specified")
	ErrTestWrongFileAmount       = errors.New("recipe rendered different amount of files than expected")
	ErrTestMissingFile           = errors.New("recipe did not render file which was expected")
	ErrTestContentMismatch       = errors.New("the contents of the files did not match")
)

func (t *Test) Validate() error {
	// TODO
	return nil
}

func (re *Recipe) RunTests() error {
	if len(re.tests) == 0 {
		return ErrNoTestsSpecified
	}

	for _, t := range re.tests {
		re.Values = t.Values
		err := re.Render()
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		if len(t.Files) != len(re.Files) {
			// TODO: show which files were missing/extra
			return ErrTestWrongFileAmount
		}

		for key, tFile := range t.Files {
			if file, ok := re.Files[key]; !ok {
				return fmt.Errorf("%w: file '%s'", ErrTestMissingFile, key)
			} else {
				if !bytes.Equal(tFile, file.Content) {
					return fmt.Errorf("%w: file '%s'.\nExpected:\n%s\n\nActual:\n%s", ErrTestContentMismatch, key, tFile, file.Content)
				}
			}

		}
	}

	return nil
}
