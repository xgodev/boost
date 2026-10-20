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
	App     string              // app name
	Package string              // package name, the app name without hyphens
	Imports []inject.ImportData // imports besides sync and fx
	Init    []string            // statements of init()
	Var     string              // prefix of the once variable
	Func    string              // prefix of the module function: <Func>Module()
	Doc     string              // rest of the doc sentence of the module function
	Name    string              // fx module name
	Options []string            // options of the fx module, as Go expressions
	// Adapter is the config switch of a @Function with a list of adapters;
	// nil for every other door.
	Adapter *adapterSwitch
}

// adapterSwitch is the config key that chooses the adapter of a @Function
// among several, and the fx options of each one.
type adapterSwitch struct {
	Key    string // config key, <app in camelCase>.function.adapter
	Names  string // the adapters, comma separated, in annotation order
	Config string // import alias of wrapper/config
	Errors string // import alias of model/errors
	Cases  []adapterCase
}

// adapterCase is one adapter of an adapterSwitch.
type adapterCase struct {
	Name    string
	Options []string // fx options of the adapter, as Go expressions
}

// modulesData feeds the template of modules.go.
type modulesData struct {
	Header  string
	App     string // app name
	Package string // package name, the app name without hyphens
	Servers bool
	Imports []inject.ImportData // the decorators, and multiserverfx when the app has servers
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
