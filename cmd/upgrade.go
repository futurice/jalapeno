package main

import (
	"bytes"
	"path/filepath"

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

	re, err := recipe.Load(source)
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	prevRe, err := recipe.LoadRendered(target, re.Name)
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	if !prevRe.IsExecuted() {
		cmd.PrintErrln("the first argument should point to the project which uses the recipe")
		return
	}

	if re.IsExecuted() {
		cmd.PrintErrln("the second argument should point to the recipe which will be used for upgrading")
		return
	}

	if re.Metadata.Name != prevRe.Metadata.Name {
		cmd.PrintErrln("recipe name used in the project should match the recipe which is used for upgrading")
		return
	}

	if semver.Compare(re.Metadata.Version, prevRe.Metadata.Version) <= 0 {
		cmd.PrintErrln("new recipe version is lower or same than the existing one")
		return
	}

	cmd.Printf("Upgrade from %s to %s\n", prevRe.Metadata.Version, re.Metadata.Version)

	re.Values = prevRe.Values

	// Check if the new version of the recipe has removed some variables
	// which existed on previous version
	for _, v := range re.Variables {
		if _, exists := re.Values[v.Name]; !exists {
			delete(re.Values, v.Name)
		}
	}

	err = recipeutil.PromptUserForValues(re)
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	err = re.Render(engine.Engine{})
	if err != nil {
		cmd.PrintErrln(err)
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
				cmd.Println("Some of the files has been manually modified. Do you want to override the following files:")
				overrideNoticed = true
			}

			// TODO: We could do better in terms of merge conflict management. Like show the diff or something
			var override bool
			prompt := &survey.Confirm{
				Message: name,
				Default: true,
			}
			err = survey.AskOne(prompt, &override)
			if err != nil {
				cmd.Println(err)
				return
			}

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
		cmd.PrintErrln(err)
		return
	}

	err = re.Save(filepath.Join(target, recipe.RenderedRecipeDirName))
	if err != nil {
		cmd.PrintErrln(err)
		return
	}
}
