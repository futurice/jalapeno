package recipe

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-yaml/yaml"
)

// Saves recipe file to given destination
func (re *Recipe) Save(dest string) error {
	out, err := yaml.Marshal(re)
	if err != nil {
		return err
	}

	// find our stack order ()
	matches, err := filepath.Glob(filepath.Join(dest, fmt.Sprintf("*-%s.yml", re.Name)))
	if err != nil {
		// The only case for this should be a malformed glob pattern
		return err
	}
	if len(matches) > 1 {
		return fmt.Errorf("this should never happen: more than one rendered %s recipe in %s", re.Name, dest)
	}

	if len(matches) == 0 {
		// No matches -- this is the first time rendering this recipe in this directory.
		// The index therefore needs to be the next highest integer, which is the length
		// of the list of all rendered recipes, plus 1.
		allRendered, err := filepath.Glob(filepath.Join(dest, "*-*.yml"))
		if err != nil {
			return err
		}
		err = os.WriteFile(filepath.Join(dest, fmt.Sprintf("%d-%s.yml", len(allRendered)+1, re.Name)), out, 0644)
		if err != nil {
			return err
		}
	} else {
		// One match -- this is the case where we're overwriting an existing version.
		err = os.WriteFile(matches[0], out, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
