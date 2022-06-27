package engine

import (
	"text/template"

	"github.com/Masterminds/sprig"
)

func funcMap() template.FuncMap {
	f := sprig.TxtFuncMap()

	// Add additional template functions here

	return f
}
