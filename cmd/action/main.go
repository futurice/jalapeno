package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/buildkite/shellwords"
	"github.com/gofrs/uuid"

	"github.com/futurice/jalapeno/internal/cli"
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
	defer output.Close() //nolint:errcheck

	// Write contents to the output file and to stdout
	out := io.MultiWriter(os.Stdout, output)

	// Since arguments are passed as a single string, we need to split them
	args, err := shellwords.SplitPosix(os.Args[1])
	checkErr(err)

	delimiter := uuid.Must(uuid.NewV4()).String()
	fmt.Fprintf(output, "%s<<%s\n", OutputResult, delimiter) //nolint:errcheck

	rootCmd := cli.NewRootCmd()
	rootCmd.SetOut(out)
	rootCmd.SetErr(out)
	rootCmd.SetArgs(args)

	exitCode := cli.Execute(rootCmd)

	fmt.Fprintf(output, "%s\n", delimiter)                   //nolint:errcheck
	fmt.Fprintf(output, "%s=%d\n", OutputExitCode, exitCode) //nolint:errcheck

	// Write buffer to the file
	err = output.Sync()
	checkErr(err)

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
