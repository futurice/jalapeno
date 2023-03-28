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

var ErrNoTestsSpecified error = errors.New("no tests specified")

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
			return fmt.Errorf("recipe rendered different amount of files than expected")
		}

		for key, tFile := range t.Files {
			if file, ok := re.Files[key]; !ok {
				return fmt.Errorf("recipe does not include file '%s' which was expected in the test case", key)
			} else {
				if !bytes.Equal(tFile, file.Content) {
					return fmt.Errorf("contents for file '%s' did not match.\nExpected:\n%s\n\nActual:\n%s", key, tFile, file.Content)
				}
			}

		}
	}

	return nil
}
