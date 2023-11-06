package main

import (
	"fmt"
	"os"

	"github.com/futurice/jalapeno/internal/cli"
)

var (
	// https://goreleaser.com/cookbooks/using-main.version/
	version string
)

func main() {
	cmd, err := cli.NewRootCmd(version)
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}

	if err = cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
