package builder

import (
	_ "embed"
	"text/template"

	"github.com/xgodev/boost/inject"
	"github.com/xgodev/boost/model/errors"
)

//go:embed templates/door.tmpl
var doorTemplate string

// doorData feeds the template of one door file (<server>_module.go or
// function_module.go).
type doorData struct {
	Header  string
	Package string              // app name
	Imports []inject.ImportData // imports besides sync and fx
	Init    []string            // statements of init()
	Var     string              // prefix of the once variable
	Func    string              // prefix of the module function: <Func>Module()
	Doc     string              // rest of the doc sentence of the module function
	Name    string              // fx module name
	Options []string            // options of the fx module, as Go expressions
}

// modulesData feeds the template of modules.go.
type modulesData struct {
	Header  string
	Package string
	Servers bool
	Imports []inject.ImportData // multiserverfx when the app has servers
	Modules []string            // module calls, in order
}

// mainData feeds the template of cmd/<app>/main.go.
type mainData struct {
	Header string
	App    string
	Import string // import path of gen/builder/<app>
}

func newTemplate() (*template.Template, error) {
	tmpl, err := template.New("builder").Parse(doorTemplate)
	if err != nil {
		return nil, errors.Wrap(err, errors.Internalf("error parsing template"))
	}
	return tmpl, nil
}
