package recipeutil

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"

	"github.com/futurice/jalapeno/pkg/recipe"
)

func IsFileModified(projectDir, path string, file recipe.File) (bool, error) {
	content, err := os.ReadFile(filepath.Join(projectDir, path))
	if err != nil {
		return false, fmt.Errorf("Could not read %s in %s: %e", path, projectDir, err)
	}
	sum := sha256.Sum256(content)
	return file.Checksum != fmt.Sprintf("sha256:%x", sum), nil
}
