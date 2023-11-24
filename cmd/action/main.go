package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/buildkite/shellwords"
	"github.com/futurice/jalapeno/internal/cli"
	"github.com/gofrs/uuid"
)

const (
	OutputResult   = "result"
	OutputExitCode = "exitcode"
)

// This is the entrypoint for the Github Action
func main() {
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		checkErr(errors.New("this image only works on Github Actions"))
	}

	output, err := os.OpenFile(os.Getenv("GITHUB_OUTPUT"), os.O_APPEND|os.O_WRONLY, 0644)
	checkErr(err)

	// Write contents to the output file and to stdout
	out := io.MultiWriter(os.Stdout, output)

	// Since arguments are passed as a single string, we need to split them
	args, err := shellwords.Split(os.Args[1])
	checkErr(err)

	delimiter := uuid.Must(uuid.NewV4()).String()
	fmt.Fprintf(output, "%s<<%s\n", OutputResult, delimiter)

	cmd := cli.NewRootCmd("")
	cmd.SetOut(out)
	cmd.SetArgs(args)

	err = cmd.ExecuteContext(context.Background())
	fmt.Fprintf(output, "%s\n", delimiter)

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
