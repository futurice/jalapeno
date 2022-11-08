package recipe

import (
	"reflect"
	"testing"
)

func TestRenderedRecipeFilesToRecipeNames(t *testing.T) {
	paths := []string{
		"foo/.jalapeno/1-bar.yml",
		"foo/.jalapeno/2-quux.yml",
	}

	recipeNames, err := renderedRecipeFilesToRecipeNames(paths)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(recipeNames, []string{"bar", "quux"}) {
		t.Errorf("Expected 'bar' and 'quux', got %s", recipeNames)
	}
}
