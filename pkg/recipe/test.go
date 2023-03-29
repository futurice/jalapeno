package recipe

import (
	"bytes"
	b64 "encoding/base64"
	"errors"
	"fmt"
)

type Test struct {
	Name   string         `yaml:"name,omitempty"`
	Values VariableValues `yaml:"values"`
	Files  TestFiles      `yaml:"files"`
}

type TestFiles map[string][]byte

func (f *TestFiles) UnmarshalYAML(unmarshal func(interface{}) error) error {
	yamlSequence := make(map[string]string)
	err := unmarshal(&yamlSequence)
	if err != nil {
		return err
	}

	files := TestFiles{}
	for name, base64Content := range yamlSequence {
		content, err := b64.StdEncoding.DecodeString(base64Content)
		if err != nil {
			return err
		}
		files[name] = content
	}

	*f = files
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
