package recipeutil_test

import (
	"testing"

	"github.com/futurice/jalapeno/pkg/recipeutil"
)

func TestExampleRecipeIsValid(t *testing.T) {
	re := recipeutil.CreateExampleRecipe("example")
	if err := re.Validate(); err != nil {
		t.Error(err)
	}
}

func TestExampleTestIsValid(t *testing.T) {
	test := recipeutil.CreateExampleTest("example")
	if err := test.Validate(); err != nil {
		t.Error(err)
	}
}
