package recipeutil

import "github.com/futurice/jalapeno/pkg/recipe"

func CreateExampleRecipe(name string) *recipe.Recipe {
	r := recipe.NewRecipe()

	variableName := "MY_VAR"
	defaultValue := "Hello World!"

	r.Metadata.Name = name
	r.Metadata.Version = "v0.0.0"
	r.Metadata.Description = "Description about what the recipe is used for and what it contains. For example tech stack, cloud environments, tools"
	r.Variables = []recipe.Variable{
		{Name: variableName, Default: defaultValue},
	}
	r.Templates = map[string][]byte{
		"README.md": []byte("{{ .Variables.MY_VAR }}"),
	}
	r.Tests = []recipe.Test{
		{
			Name:   "defaults",
			Values: recipe.VariableValues{variableName: defaultValue},
			Files: map[string][]byte{
				"README.md": []byte(defaultValue),
			},
		},
	}

	return r
}

func CreateExampleTest() *recipe.Test {
	test := &recipe.Test{
		Name: "example",
	}

	return test
}
