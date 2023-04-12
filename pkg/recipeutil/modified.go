package recipeutil

import (
	"crypto/sha256"
	"fmt"
	"os"

	"github.com/futurice/jalapeno/pkg/recipe"
)

func IsFileModified(path string, file recipe.File) (bool, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return false, fmt.Errorf("could not read %s: %e", path, err)
	}
	sum := sha256.Sum256(content)
	return file.Checksum != fmt.Sprintf("sha256:%x", sum), nil
}
