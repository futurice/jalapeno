package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/recipeutil"
)

// Utility function for creating a retry message for the user such that they can re-run the cli command with the same values
func makeRetryMessage(values recipe.VariableValues) string {
	var commandline strings.Builder
	skipNext := false
	for idx, arg := range os.Args {
		if skipNext {
			skipNext = false
			continue
		}
		if arg == "--set" {
			skipNext = true
			continue
		}
		// quote all non-option args to be on the safe side, except for indices 0 and 1
		// (which are the program name and the command name)
		if idx <= 1 || strings.HasPrefix(arg, "-") {
			commandline.WriteString(arg)
		} else {
			commandline.WriteString(fmt.Sprintf("\"%s\"", arg))
		}
		commandline.WriteString(" ")
	}

	for key, value := range values {
		commandline.WriteString(" --set ")
		switch value := value.(type) {
		case []map[string]string: // serialize to CSV
			csv, err := recipeutil.TableToCSV(value, ',')
			if err != nil {
				panic(err)
			}
			commandline.WriteString(fmt.Sprintf("\"%s=%s\"", key, strings.ReplaceAll(strings.TrimRight(csv, "\n"), "\n", "\\n")))
		default:
			commandline.WriteString(fmt.Sprintf("\"%s=%s\"", key, value))
		}
	}
	retryMessage := fmt.Sprintf("To re-run the recipe with the same values, use the following command:\n\n%s", commandline.String())
	return retryMessage
}
