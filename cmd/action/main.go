package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/futurice/jalapeno/internal/cli"
)

const (
	OutputExitCode = "exitcode"
)

// This is the entrypoint for the Github Action
func main() {
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		checkErr(errors.New("this image only works on Github Actions"))
	}

	output, err := os.OpenFile(os.Getenv("GITHUB_OUTPUT"), os.O_APPEND|os.O_WRONLY, 0644)
	checkErr(err)

	cmd := cli.NewRootCmd("")
	err = cmd.ExecuteContext(context.Background())

	exitCode, isExitCodeSet := cmd.Context().Value(cli.ExitCodeContextKey{}).(int)
	if !isExitCodeSet {
		if err == nil {
			exitCode = 0
		} else {
			exitCode = 1
		}
	}
	fmt.Fprintf(output, "%s=%d\n", OutputExitCode, exitCode)

	// Write buffer to the file
	err = output.Sync()
	checkErr(err)

	output.Close()

	// Map all non error exit codes to 0 so that Github Actions job does not fail
	if exitCode != cli.ExitCodeError {
		os.Exit(cli.ExitCodeOK)
	} else {
		os.Exit(cli.ExitCodeError)
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(cli.ExitCodeError)
	}
}
