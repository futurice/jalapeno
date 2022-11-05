package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/futurice/jalapeno/pkg/recipe"
	"github.com/spf13/cobra"
)

var (
	targetPath = ""
)

func newEjectCmd() *cobra.Command {
	var ejectCmd = &cobra.Command{
		Use:   "eject",
		Short: "Remove all Jalapeno-specific files",
		Long:  "Remove all the files and directories that are for Jalapeno internal use, and leave only the rendered project files.",
		Args:  cobra.ExactArgs(0),
		Run:   ejectFunc,
	}

	ejectCmd.Flags().StringVarP(&targetPath, "path", "p", ".", "Location of project to eject")

	return ejectCmd
}

func ejectFunc(cmd *cobra.Command, args []string) {
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		fmt.Println("Project path does not exist")
		return
	}

	jalapenoPath := filepath.Join(targetPath, recipe.RenderedRecipeDirName)

	if stat, err := os.Stat(jalapenoPath); os.IsNotExist(err) || !stat.IsDir() {
		fmt.Printf("'%s' is not a Jalapeno project\n", targetPath)
		return
	}

	fmt.Printf("Deleting %s...", jalapenoPath)
	err := os.RemoveAll(jalapenoPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("\nEjected successfully!")
}
