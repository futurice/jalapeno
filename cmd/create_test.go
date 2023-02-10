package main

import "testing"

func TestExampleRecipe(t *testing.T) {
	recipe := createExampleRecipe("foo")
	if err := recipe.Validate(); err != nil {
		t.Errorf("failed to validate the example recipe: %s", err)
	}
}
