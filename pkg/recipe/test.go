package recipe

import (
	"bytes"
	b64 "encoding/base64"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
)

type Test struct {
	// Name of the test case. Defined by the test filename
	Name string `yaml:"-"`

	// Values to use to render the recipe templates
	Values VariableValues `yaml:"values"`

	// Snapshots of the rendered templates which were rendered with the values specified in the test
	Files map[string]TestFile `yaml:"files"`
}

type TestFile []byte

// Random hardcoded UUID
var TestID uuid.UUID = uuid.Must(uuid.FromString("12345678-1234-5678-1234-567812345678"))

func (f TestFile) MarshalYAML() (interface{}, error) {
	return b64.StdEncoding.EncodeToString(f), nil
}

func (f *TestFile) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var base64String string
	err := unmarshal(&base64String)
	if err != nil {
		return err
	}

	file, err := b64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return err
	}

	*f = file
	return nil
}

var (
	ErrNoTestsSpecified    error = errors.New("no tests specified")
	ErrTestWrongFileAmount       = errors.New("recipe rendered different amount of files than expected")
	ErrTestMissingFile           = errors.New("recipe did not render file which was expected")
	ErrTestContentMismatch       = errors.New("the contents of the files did not match")
)

func (t *Test) Validate() error {
	if t.Name == "" {
		return errors.New("test name can not be empty")
	}

	return nil
}

func (re *Recipe) RunTests() []error {
	errors := make([]error, len(re.Tests))
	for i, t := range re.Tests {
		sauce, err := re.Execute(t.Values, TestID)
		if err != nil {
			errors[i] = fmt.Errorf("%w", err)
			continue
		}

		if len(t.Files) != len(sauce.Files) {
			// TODO: show which files were missing/extra
			errors[i] = ErrTestWrongFileAmount
			continue
		}

		for key, tFile := range t.Files {
			if file, ok := sauce.Files[key]; !ok {
				errors[i] = fmt.Errorf("%w: file '%s'", ErrTestMissingFile, key)
				continue
			} else {
				if !bytes.Equal(tFile, file.Content) {
					errors[i] = fmt.Errorf("%w: file '%s'.\nExpected:\n%s\n\nActual:\n%s", ErrTestContentMismatch, key, tFile, file.Content)
					continue
				}
			}
		}
	}

	return errors
}
