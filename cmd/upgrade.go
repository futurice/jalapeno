package main

import (
	"bytes"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/futurice/jalapeno/pkg/engine"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/recipeutil"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

func newUpgradeCmd() *cobra.Command {
	// upgradeCmd represents the upgrade command
	var upgradeCmd = &cobra.Command{
		Use:   "upgrade PROJECT RECIPE",
		Short: "Upgrade recipe in a project",
		Long:  "", // TODO
		Run:   upgradeFunc,
		Args:  cobra.ExactArgs(2),
	}

	return upgradeCmd
}

func upgradeFunc(cmd *cobra.Command, args []string) {
	target := args[0]
	source := args[1]

	prevRe, err := recipe.LoadRenderedFromDir(target)
	if err != nil {
		fmt.Println(err)
		return
	}

	re, err := recipe.LoadFromDir(source)
	if err != nil {
		fmt.Println(err)
		return
	}

	if semver.Compare(re.Metadata.Version, prevRe.Metadata.Version) <= 0 {
		fmt.Println("error: new recipe version is lower or same than the existing one")
		return
	}

	fmt.Printf("Upgrade from %s to %s\n", re.Metadata.Version, prevRe.Metadata.Version)

	re.Values = prevRe.Values

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
	overrideNoticed := false

	for name := range re.Files {
		if _, exists := prevRe.Files[name]; exists {
			if bytes.Equal(re.Files[name], prevRe.Files[name]) {
				// A file with exactly same name and content already exist, skip
				continue
			}

			// The file contents has been modified

			if !overrideNoticed {
				fmt.Println("Some of the files has been manually modified. Do you want to override the following files:")
				overrideNoticed = true
			}

			// TODO: We could do better in terms of merge conflict management. Like show the diff or something
			var override bool
			prompt := &survey.Confirm{
				Message: name,
				Default: true,
			}
			survey.AskOne(prompt, &override)
			if !override {
				// User decided not to override the file with manual changes, we can skip the file
				continue
			}
		}

		// Add new file or replace existing one
		output[name] = re.Files[name]
	}

	err = recipeutil.SaveFiles(output, target)
	if err != nil {
		fmt.Println(err)
		return
	}
}
