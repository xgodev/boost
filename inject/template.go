package inject

import (
	_ "embed"
	"text/template"

	"github.com/xgodev/boost/model/errors"
)

//go:embed templates/module.tmpl
var moduleTemplate string

// ModuleData feeds the module template of one component.
type ModuleData struct {
	Package string       // package name of the generated file
	Func    string       // annotated function name
	Source  string       // import path of the annotated function
	Kind    string       // "Provide" or "Invoke"
	Target  string       // expression passed to fx.Provide/fx.Invoke
	Deps    []string     // module calls of the components it depends on
	Imports []ImportData // imports besides sync and fx
}

// RootData feeds the template of the root module.go.
type RootData struct {
	Package string
	Imports []ImportData
	Modules []string
}

// ImportData is an aliased import of a generated file.
type ImportData struct {
	Alias string
	Path  string
}

// NewTemplate parses the generator templates ("module" and "root").
func NewTemplate() (*template.Template, error) {
	tmpl, err := template.New("inject").Parse(moduleTemplate)
	if err != nil {
		return nil, errors.Wrap(err, errors.Internalf("error parsing template"))
	}
	return tmpl, nil
}
