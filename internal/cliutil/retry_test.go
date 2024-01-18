package cliutil_test

import (
	"strings"
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
			`jalapeno execute "path/to/recipe"`,
		},
		{
			"Multiple string values",
			[]string{"jalapeno", "execute", "path/to/recipe"},
			recipe.VariableValues{
				"key1": "value1",
				"key2": "value2",
			},
			`jalapeno execute "path/to/recipe" --set "key1=value1" --set "key2=value2"`,
		},
		{
			"Boolean values",
			[]string{"jalapeno", "execute", "path/to/recipe"},
			recipe.VariableValues{
				"key1": true,
				"key2": false,
			},
			`jalapeno execute "path/to/recipe" --set "key1=true" --set "key2=false"`,
		},
		{
			"Table values",
			[]string{"jalapeno", "execute", "path/to/recipe"},
			recipe.VariableValues{
				"table1": recipe.TableValue{
					Columns: []string{"col1", "col2"},
					Rows: [][]string{
						{"value1", "value2"},
					},
				},
				"table2": recipe.TableValue{
					Columns: []string{"col3", "col4"},
					Rows: [][]string{
						{"value3", "value4"},
					},
				},
			},
			`jalapeno execute "path/to/recipe" --set "table1=value1,value2" --set "table2=value3,value4"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(tt *testing.T) {
			result := cliutil.MakeRetryMessage(tc.Args, tc.Values)
			if !strings.Contains(result, tc.Expected) {
				tt.Errorf("unexpected result:\n%s", diff.Diff(tc.Expected, result))
			}
		})
	}
}
