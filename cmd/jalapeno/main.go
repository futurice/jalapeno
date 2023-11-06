package main

import (
	"context"
	"os"

	"github.com/futurice/jalapeno/internal/cli"
)

var (
	// https://goreleaser.com/cookbooks/using-main.version/
	version string
)

func main() {
	cmd := cli.NewRootCmd(version)
	err := cmd.ExecuteContext(context.Background())

	exitCode := cmd.Context().Value(cli.ExitCodeContextKey{})
	if code, ok := exitCode.(int); ok { // Make sure that the exit code is still an int
		os.Exit(code) // Exit with the exit code defined by a subcommand
	} else {
		if err == nil {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}
