package recipe

import (
	"crypto/sha256"
	"fmt"
	"path/filepath"
)

type File struct {
	Checksum      string `yaml:"checksum"` // e.g. "sha256:asdjfajdfa" w. default algo
	Content       []byte `yaml:"-"`
	IgnoreUpgrade bool   `yaml:"-"` // default value is false from golang bool default
}

type Recipe struct {
	Metadata  `yaml:",inline"`
	Variables []Variable        `yaml:"vars,omitempty"`
	Values    VariableValues    `yaml:"values,omitempty"`
	Templates map[string][]byte `yaml:"-"`
	Files     map[string]File   `yaml:"files,omitempty"`
}

type RenderEngine interface {
	Render(recipe *Recipe, values map[string]interface{}) (map[string][]byte, error)
}

func (re *Recipe) Validate() error {
	if err := re.Metadata.Validate(); err != nil {
		return err
	}

	checkDuplicates := make(map[string]bool)
	for _, v := range re.Variables {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("error on variable %s: %w", v.Name, err)
		}
		if _, exists := checkDuplicates[v.Name]; exists {
			return fmt.Errorf("variable %s has been declared multiple times", v.Name)
		}
		checkDuplicates[v.Name] = true
	}

	return nil
}

// Renders recipe templates from .Templates to .Files
func (re *Recipe) Render(engine RenderEngine) error {
	// Define the context which is available on templates
	context := map[string]interface{}{
		"Recipe":    re.Metadata,
		"Variables": re.Values,
	}

	var err error
	files, err := engine.Render(re, context)
	if err != nil {
		return err
	}

	re.Files = make(map[string]File, len(files))
	idx := 0
	for filename, content := range files {
		sum := sha256.Sum256(content)
		file := File{Content: content, Checksum: fmt.Sprintf("sha256:%x", sum)}
		for _, pattern := range re.Metadata.IgnorePatterns {
			if matched, err := filepath.Match(pattern, filename); err != nil {
				return fmt.Errorf("bad ignore pattern %s: %w", pattern, err)
			} else if matched {
				// mark file as ignored
				file.IgnoreUpgrade = true
			}
		}
		re.Files[filename] = file
		idx += 1
		if idx > len(files) {
			return fmt.Errorf("Files array grew during execution")
		}
	}

	return nil
}

// Check if the recipe is in executed state (the templates has been rendered)
func (re *Recipe) IsExecuted() bool {
	return len(re.Files) > 0
}

type RecipeConflict struct {
	Path           string
	Sha256Sum      string
	OtherSha256Sum string
}

// Check if the recipe conflicts with another recipe. Recipes conflict if they touch the same files.
func (re *Recipe) Conflicts(other *Recipe) []RecipeConflict {
	var conflicts []RecipeConflict
	for path, file := range re.Files {
		if otherFile, exists := other.Files[path]; exists {
			conflicts = append(
				conflicts,
				RecipeConflict{
					Path:           path,
					Sha256Sum:      file.Checksum,
					OtherSha256Sum: otherFile.Checksum,
				})
		}
	}
	return conflicts
}
