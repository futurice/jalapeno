package main

import (
	"fmt"
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
		Use:     "execute <recipe_path>",
		Aliases: []string{"exec", "e"},
		Short:   "Execute a given recipe and save output to path",
		Args:    cobra.ExactArgs(1),
		Run:     executeFunc,
	}

	execCmd.Flags().StringVarP(&outputBasePath, "output", "o", ".", "Path where the output files should be created")

	return execCmd
}

func executeFunc(cmd *cobra.Command, args []string) {
	if _, err := os.Stat(outputBasePath); os.IsNotExist(err) {
		fmt.Println("Output path does not exist")
		return
	}

	re, err := recipe.Load(args[0])
	if err != nil {
		fmt.Printf("Error when loading the recipe: %s\n", err)
		return
	}

	fmt.Printf("Recipe name: %s\n\n", re.Metadata.Name)

	err = re.Validate()
	if err != nil {
		fmt.Printf("The provided recipe was invalid: %s\n", err)
		return
	}

	// TODO: Set values provided by --set flag to re.Values

	err = recipeutil.PromptUserForValues(re)
	if err != nil {
		fmt.Printf("Error when prompting for values: %s\n", err)
		return
	}

	err = re.Render(engine.Engine{})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create sub directory for recipe
	recipePath := filepath.Join(outputBasePath, recipe.RenderedRecipeDirName)
	err = os.MkdirAll(recipePath, 0700)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = re.Save(recipePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = recipeutil.SaveFiles(re.RenderedTemplates, outputBasePath)
	if err != nil {
		fmt.Println(err)
		return
	}
}
