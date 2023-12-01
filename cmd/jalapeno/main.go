package main

import (
	"os"

	"github.com/futurice/jalapeno/internal/cli"
)

// This is the entrypoint for the CLI
func main() {
	rootCmd := cli.NewRootCmd()
	exitCode := cli.Execute(rootCmd)
	os.Exit(exitCode)
}
