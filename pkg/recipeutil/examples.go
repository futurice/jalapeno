package recipeutil

import "github.com/futurice/jalapeno/pkg/recipe"

func CreateExampleRecipe(name string) recipe.Recipe {
	r := recipe.NewRecipe()

	variableName := "MY_VAR"
	defaultValue := "Hello World!"

	r.Metadata.Name = name
	r.Metadata.Version = "v0.0.0"
	r.Metadata.Description = "Description about what the recipe is used for and what it contains. For example tech stack, cloud environments, tools"
	r.Variables = []recipe.Variable{
		{Name: variableName, Default: defaultValue},
	}
	r.Templates = map[string]recipe.File{
		"README.md": recipe.NewFile([]byte("{{ .Variables.MY_VAR }}")),
	}
	r.Tests = []recipe.Test{
		{
			Name:   "defaults",
			Values: recipe.VariableValues{variableName: defaultValue},
			Files: map[string]recipe.File{
				"README.md": recipe.NewFile([]byte(defaultValue)),
			},
		},
	}

	return r
}

func CreateExampleTest(name string) recipe.Test {
	return recipe.Test{
		Name: name,
	}
}

func CreateExampleManifest() recipe.Manifest {
	m := recipe.NewManifest()

	m.Recipes = []recipe.ManifestRecipe{
		{
			Name:       "recipe-a",
			Version:    "v0.0.1",
			Repository: "./path/to/recipe-a",
		},
		{
			Name:       "recipe-b",
			Version:    "v0.0.1",
			Repository: "oci://url/to/recipe-b",
			Values: recipe.VariableValues{
				"MY_VAR": "Hello World!",
			},
		},
	}

	return m
}
