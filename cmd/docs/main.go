package main

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"text/template"

	"github.com/futurice/jalapeno/internal/cli"
	"github.com/spf13/pflag"
)

type Flag struct {
	Name        string
	Shorthand   string
	Default     string
	Description string
}

type CommandInfo struct {
	Name        string
	Description string
	Usage       string
	Flags       []Flag
}

//go:embed templates
var tmpls embed.FS

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Error: no destination path provided")
		return
	}

	rootCmd, err := cli.NewRootCmd()
	checkErr(err)

	cmds := rootCmd.Commands()
	infos := make([]CommandInfo, len(cmds))
	for i, c := range cmds {
		info := CommandInfo{
			Name:        c.Name(),
			Description: c.Long,
			Usage:       c.Use,
			Flags:       make([]Flag, 0),
		}

		c.Flags().VisitAll(func(f *pflag.Flag) {
			info.Flags = append(info.Flags, Flag{
				Name:        f.Name,
				Shorthand:   f.Shorthand,
				Default:     f.DefValue,
				Description: f.Usage,
			})
		})

		infos[i] = info
	}

	tmpl := template.Must(template.New("doc").ParseFS(tmpls, "templates/*"))

	var b bytes.Buffer
	err = tmpl.ExecuteTemplate(&b, "main.tmpl", infos)
	checkErr(err)

	err = os.WriteFile(args[0], b.Bytes(), 0644)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}
}
