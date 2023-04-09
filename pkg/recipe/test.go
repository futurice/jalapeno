package recipe

import (
	"bytes"
	b64 "encoding/base64"
	"errors"
	"fmt"
)

type Test struct {
	Name   string              `yaml:"name,omitempty"`
	Values VariableValues      `yaml:"values"`
	Files  map[string]TestFile `yaml:"files"`
}

type TestFile []byte

func (f *TestFile) MarshalYAML() (interface{}, error) {
	return b64.StdEncoding.EncodeToString(*f), nil
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
	// TODO
	return nil
}

func (re *Recipe) RunTests() []error {
	errors := make([]error, len(re.Tests))
	for i, t := range re.Tests {
		sauce, err := re.Execute(t.Values)
		if err != nil {
			errors[i] = fmt.Errorf("%w", err)
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
