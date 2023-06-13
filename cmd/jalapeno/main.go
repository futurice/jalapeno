package main

import (
	"fmt"
	"os"

	"github.com/futurice/jalapeno/internal/cli"
)

func main() {
	cmd, err := cli.NewRootCmd(os.Stdout, os.Args[1:])
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}

	if err = cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
