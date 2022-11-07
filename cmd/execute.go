package main

import (
	"os"
	"path/filepath"

	"github.com/futurice/jalapeno/pkg/engine"
	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/futurice/jalapeno/pkg/recipeutil"
	"github.com/spf13/cobra"
)

var (
	outputBasePath = ""
)

func newExecuteCmd() *cobra.Command {
	var execCmd = &cobra.Command{
		Use:     "execute RECIPE",
		Aliases: []string{"exec", "e"},
		Short:   "Execute a given recipe and save output to path",
		Long:    "", // TODO
		Args:    cobra.ExactArgs(1),
		Run:     executeFunc,
	}

	execCmd.Flags().StringVarP(&outputBasePath, "output", "o", ".", "Path where the output files should be created")

	return execCmd
}

func executeFunc(cmd *cobra.Command, args []string) {
	if _, err := os.Stat(outputBasePath); os.IsNotExist(err) {
		cmd.PrintErrln("Error: output path does not exist")
		return
	}

	re, err := recipe.Load(args[0])
	if err != nil {
		cmd.PrintErrf("Error: can't load the recipe: %s\n", err)
		return
	}

	cmd.Printf("Recipe name: %s\n", re.Metadata.Name)

	if re.Metadata.Description != "" {
		cmd.Printf("Description: %s\n", re.Metadata.Description)
	}

	err = re.Validate()
	if err != nil {
		cmd.PrintErrf("Error: the provided recipe was invalid: %s\n", err)
		return
	}

	if len(re.Templates) == 0 {
		cmd.PrintErrf("Error: the recipe does not contain any templates")
		return
	}

	// TODO: Set values provided by --set flag to re.Values

	err = recipeutil.PromptUserForValues(re)
	if err != nil {
		cmd.PrintErrf("Error when prompting for values: %s\n", err)
		return
	}

	err = re.Render(engine.Engine{})
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	// Create sub directory for recipe
	recipePath := filepath.Join(outputBasePath, recipe.RenderedRecipeDirName)
	err = os.MkdirAll(recipePath, 0700)
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	err = re.Save(recipePath)
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	err = recipeutil.SaveFiles(re.Files, outputBasePath)
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	cmd.Println("\nRecipe executed successfully!")

	if re.InitHelp != "" {
		cmd.Printf("Next up: %s\n", re.InitHelp)
	}
}
