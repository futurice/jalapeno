package cli_test

import (
	"testing"

	"github.com/futurice/jalapeno/internal/cli"
)

func TestExampleRecipe(t *testing.T) {
	recipe := cli.CreateExampleRecipe("foo")
	if err := recipe.Validate(); err != nil {
		t.Fatalf("failed to validate the example recipe: %s", err)
	}
}
