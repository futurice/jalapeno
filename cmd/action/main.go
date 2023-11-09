package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/futurice/jalapeno/internal/cli"
	"github.com/gofrs/uuid"
)

var (
	// https://goreleaser.com/cookbooks/using-main.version/
	version string
)

func main() {
	filename := os.Getenv("GITHUB_OUTPUT")
	if filename == "" {
		checkErr(errors.New("GITHUB_OUTPUT environment variable not set"))
	}

	output, err := os.OpenFile(filename, os.O_APPEND, 0644)
	checkErr(err)

	cmd := cli.NewRootCmd(version)
	delimiter := uuid.Must(uuid.NewV4()).String()
	fmt.Fprintf(output, "result<<%s\n", delimiter)

	// TODO: Add outputs

	fmt.Fprintf(output, "%s\n", delimiter)

	err = cmd.Execute()
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}
