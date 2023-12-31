package cliutil_test

import (
	"testing"

	"github.com/futurice/jalapeno/internal/cliutil"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/kylelemons/godebug/diff"
)

func TestMakeRetryMessage(t *testing.T) {
	testCases := []struct {
		Name     string
		Args     []string
		Values   recipe.VariableValues
		Expected string
	}{
		{
			"No values",
			[]string{"jalapeno", "execute", "path/to/recipe"},
			recipe.VariableValues{},
			`To re-run the recipe with the same values, use the following command:

jalapeno execute "path/to/recipe"`,
		},
		{
			"Non-empty values",
			[]string{"jalapeno", "execute", "path/to/recipe"},
			recipe.VariableValues{
				"key1": "value1",
				"key2": "value2",
			},
			`To re-run the recipe with the same values, use the following command:

jalapeno execute "path/to/recipe" --set "key1=value1" --set "key2=value2"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(tt *testing.T) {
			result := cliutil.MakeRetryMessage(tc.Args, tc.Values)
			if result != tc.Expected {
				tt.Errorf("unexpected result:\n%s", diff.Diff(tc.Expected, result))
			}
		})
	}
}
