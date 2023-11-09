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

	exitCode, isExitCodeSet := cmd.Context().Value(cli.ExitCodeContextKey{}).(int)
	if !isExitCodeSet {
		if err == nil {
			exitCode = 0
		} else {
			exitCode = 1
		}
	}

	os.Exit(exitCode)
}
