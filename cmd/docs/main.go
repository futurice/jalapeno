package main

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"text/template"

	"github.com/futurice/jalapeno/internal/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Flag struct {
	Name        string
	Shorthand   string
	Default     string
	Description string
	Type        string
}

type CommandInfo struct {
	Name        string
	Description string
	Usage       string
	Example     string
	Flags       []Flag
}

//go:embed templates
var tmpls embed.FS

var (
	// https://goreleaser.com/cookbooks/using-main.version/
	version string
)

// This is the entrypoint for generating API reference documentation
func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Error: no destination path provided")
		return
	}

	rootCmd := cli.NewRootCmd(version)
	subCmds := rootCmd.Commands()

	tmpl := template.Must(template.New("doc").ParseFS(tmpls, "templates/*"))

	var b bytes.Buffer
	err := tmpl.ExecuteTemplate(&b, "main.tmpl", map[string]interface{}{
		"Commands": mapCommandInfos(subCmds),
	})
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

func mapCommandInfos(cmds []*cobra.Command) []CommandInfo {
	infos := make([]CommandInfo, len(cmds))
	for i, c := range cmds {
		info := CommandInfo{
			Name:        c.Name(),
			Description: c.Long,
			Usage:       c.Use,
			Example:     c.Example,
			Flags:       make([]Flag, 0),
		}

		c.Flags().VisitAll(func(f *pflag.Flag) {
			info.Flags = append(info.Flags, Flag{
				Name:        f.Name,
				Shorthand:   f.Shorthand,
				Default:     f.DefValue,
				Type:        valueTypeToString(f.Value),
				Description: f.Usage,
			})
		})

		infos[i] = info
	}

	return infos
}

func valueTypeToString(v pflag.Value) string {
	t := v.Type()
	switch t {
	case "stringArray":
		return "[]string"
	default:
		return t
	}

}
