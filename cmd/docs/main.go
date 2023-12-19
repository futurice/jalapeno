package main

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"text/template"

	"github.com/Masterminds/sprig"
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
	Aliases     []string
	Description string
	Usage       string
	Example     string
	Flags       []Flag
}

//go:embed all:templates
var tmpls embed.FS

const (
	referenceDocPath = "./docs/site/docs/api.md"
	changelogSource  = "./CHANGELOG.md"
	changelogTarget  = "./docs/site/docs/changelog.md"
)

var templates = template.Must(template.
	New("doc").
	Funcs(sprig.FuncMap()).
	ParseFS(tmpls, "templates/*"),
)

// This is the entrypoint for generating API reference documentation
func main() {
	err := GenerateReferenceDoc()
	checkErr(err)
	fmt.Println("Reference documentation generated")

	err = GenerateChangelog()
	checkErr(err)
	fmt.Println("Changelog generated")
}

func GenerateChangelog() error {
	changelog, err := os.ReadFile(changelogSource)
	if err != nil {
		return err
	}

	var b bytes.Buffer
	err = templates.ExecuteTemplate(&b, "changelog.tmpl", map[string]interface{}{
		"Changelog": string(changelog),
	})
	if err != nil {
		return err
	}

	// Write the changelog to the target file
	return os.WriteFile(changelogTarget, b.Bytes(), 0644)
}

func GenerateReferenceDoc() error {
	rootCmd := cli.NewRootCmd()
	subCmds := rootCmd.Commands()

	var b bytes.Buffer
	err := templates.ExecuteTemplate(&b, "reference.tmpl", map[string]interface{}{
		"Commands": mapCommandInfos(subCmds),
	})
	if err != nil {
		return err
	}

	err = os.WriteFile(referenceDocPath, b.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
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
			Aliases:     c.Aliases,
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
	switch t := v.Type(); t {
	case "stringArray":
		return "[]string"
	default:
		return t
	}
}
