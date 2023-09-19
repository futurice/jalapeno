package main

import (
	"fmt"
	"os"
	"time"

	"github.com/carlmjohnson/versioninfo"
	"github.com/futurice/jalapeno/internal/cli"
)

var (
	// https://goreleaser.com/cookbooks/using-main.version/
	version string
)

func main() {
	cmd, err := cli.NewRootCmd()
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}

	if version != "" {
		cmd.Version = version
	} else {
		cmd.Version = fmt.Sprintf(
			"%s (Built on %s from Git SHA %s)",
			versioninfo.Version,
			versioninfo.Revision,
			versioninfo.LastCommit.Format(time.RFC3339),
		)

	}

	if err = cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
