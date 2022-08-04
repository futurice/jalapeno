package main

import (
	"bytes"
	"fmt"

	"github.com/futurice/jalapeno/pkg/engine"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/recipeutil"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

func newUpgradeCmd() *cobra.Command {
	// upgradeCmd represents the upgrade command
	var upgradeCmd = &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade recipe in a project",
		Long:  "",
		Run:   upgradeFunc,
	}

	return upgradeCmd
}

func upgradeFunc(cmd *cobra.Command, args []string) {
	source := "./examples/gcp-web-server" // TODO
	target := "./dist"                    // TODO

	re, err := recipe.LoadFromDir(source)
	if err != nil {
		fmt.Println(err)
		return
	}

	exRe, err := recipe.LoadRenderedFromDir(target)
	if err != nil {
		fmt.Println(err)
		return
	}

	if semver.Compare(re.Metadata.Version, exRe.Metadata.Version) <= 0 {
		fmt.Println("error: new recipe version is lower or same than the existing one")
		return
	}

	fmt.Printf("Upgrade from %s to %s", re.Metadata.Version, exRe.Metadata.Version)

	re.Values = exRe.Values

	// TODO: Clean up values which does not exist in the new recipe

	err = recipeutil.PromptUserForValues(re)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = re.Render(engine.Engine{})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Collect files which should be written to the destination directory
	output := make(map[string][]byte)

	for name := range re.Files {
		if _, exists := exRe.Files[name]; exists {
			if bytes.Compare(re.Files[name], exRe.Files[name]) != 0 {
				fmt.Printf("%s: MODIFIED\n", name)
				// TODO: Apply merge, handle conflicts
			} else {
				fmt.Printf("%s: KEEP\n", name)
				continue
			}
		} else {
			fmt.Printf("%s: NEW\n", name)
			output[name] = re.Files[name]
		}
	}

	// err = recipeutil.SaveFiles(output, target)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
}
