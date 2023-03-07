package main

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/futurice/jalapeno/cmd/internal/option"
	"github.com/futurice/jalapeno/pkg/engine"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/recipeutil"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

type upgradeOptions struct {
	TargetPath string
	SourcePath string
	option.Common
}

func newUpgradeCmd() *cobra.Command {
	var opts upgradeOptions
	var cmd = &cobra.Command{
		Use:   "upgrade PROJECT RECIPE",
		Short: "Upgrade recipe in a project",
		Long:  "", // TODO
		Args:  cobra.ExactArgs(2),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.TargetPath = args[0]
			opts.SourcePath = args[1]
			return option.Parse(&opts)
		},
		Run: func(cmd *cobra.Command, args []string) {
			runUpgrade(cmd, opts)
		},
	}

	if err := option.ApplyFlags(&opts, cmd.Flags()); err != nil {
		return nil
	}

	return cmd
}

func runUpgrade(cmd *cobra.Command, opts upgradeOptions) {
	re, err := recipe.Load(opts.SourcePath)
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	rendered, err := recipe.LoadRendered(opts.TargetPath)
	if err != nil {
		cmd.PrintErrln(err)
		return
	}
	var prevRe *recipe.Recipe
	for _, r := range rendered {
		if r.Name == re.Name {
			prevRe = &r
			break
		}
	}
	if prevRe == nil {
		cmd.PrintErrf("directory %s does not contain recipe %s\n", opts.TargetPath, re.Name)
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
	}

	err = re.Render(engine.Engine{})
	if err != nil {
		return
	}

	// read common ignore file if it exists
	ignorePatterns := make([]string, 0)
	if data, err := os.ReadFile(filepath.Join(opts.TargetPath, recipe.IgnoreFileName)); err == nil {
		ignorePatterns = append(ignorePatterns, strings.Split(string(data), "\n")...)
	} else if !errors.Is(err, fs.ErrNotExist) {
		// something else happened than trying to read an ignore file that does not exist
		cmd.PrintErrf("failed to read ignore file: %v\n", err)
		return
	}
	ignorePatterns = append(ignorePatterns, re.IgnorePatterns...)

	// Collect files which should be written to the destination directory
	output := make(map[string]recipe.File, len(re.Files))
	overrideNoticed := false

	for path := range re.Files {
		skip := false
		for _, pattern := range ignorePatterns {
			if matched, err := filepath.Match(pattern, path); err != nil {
				cmd.PrintErrf("bad ignore pattern '%s': %v\n", pattern, err)
				return
			} else if matched {
				// file was marked as ignored for upgrades
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		if prevFile, exists := prevRe.Files[path]; exists {
			// Check if file was modified after rendering
			if modified, err := recipeutil.IsFileModified(opts.TargetPath, path, prevFile); err != nil {
				cmd.PrintErrln(err)
				return
			} else if modified {
				// The file contents has been modified
				if !overrideNoticed {
					cmd.Println("Some of the files has been manually modified. Do you want to override the following files:")
					overrideNoticed = true
				}

				// TODO: We could do better in terms of merge conflict management. Like show the diff or something
				var override bool
				prompt := &survey.Confirm{
					Message: path,
					Default: true,
				}
				err = survey.AskOne(prompt, &override)
				if err != nil {
					cmd.PrintErrln(err)
					return
				}

				if !override {
					// User decided not to override the file with manual changes, remove from
					// list of changes to write
					continue
				}
			}
		}

		// Add new file or replace existing one
		output[path] = re.Files[path]
	}

	err = recipeutil.SaveFiles(output, opts.TargetPath)
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	err = re.Save(opts.TargetPath)
	if err != nil {
		cmd.PrintErrln(err)
		return
	}
}
