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

// ImplementationData feeds the template of the switch of an implementation
// set.
type ImplementationData struct {
	Package string               // package name of the generated file
	Ident   string               // exported prefix: <Ident>Implementation, <Ident>ImplementationModule
	Var     string               // unexported prefix of the key constant and the once guard
	Type    string               // path.Type of the switched value, for the docs
	Key     string               // config key
	Default string               // name of the default implementation
	Names   string               // implementation names, comma separated
	Cases   []ImplementationCase // sorted by name
	Imports []ImportData         // generated packages of the implementations
}

// ConfigData feeds the template of the module of a @Config struct.
type ConfigData struct {
	Package string          // package name of the generated file
	Struct  string          // name of the @Config struct
	Var     string          // unexported prefix of the root key constant and the once guard
	Source  string          // import path of the struct's package
	Alias   string          // import alias of the struct's package
	Root    string          // root config key of the struct
	Keys    []ConfigKeyData // in field order
	Time    bool            // a default is a time.Duration
}

// ConfigKeyData is one config key of a @Config struct.
type ConfigKeyData struct {
	Segment string // key segment under the root
	Value   string // Go literal of the default
	Desc    string // Go literal of the description
	Env     string // environment variable that sets the key
	Hide    bool   // the value is a secret
}

// ClientData feeds the template of a client inject builds.
type ClientData struct {
	Package    string             // package name of the generated file
	Ident      string             // exported prefix: <Ident>Module
	Var        string             // unexported prefix of the root constant, the once guard and the constructor
	Desc       string             // the client, for the docs
	Root       string             // config root of the client
	Env        string             // environment variables that set its keys
	Module     string             // name of the fx module
	Target     string             // expression passed to fx.Provide
	Type       string             // Go type of the client
	Build      string             // ...WithConfigPath constructor of the factory
	Constraint string             // build constraint of the file, "" when it builds anywhere
	Plugins    []ClientPluginData // in the declared order
	Names      string             // plugin names, comma separated
	Close      bool               // the module closes the client when the app stops
	Imports    []ImportData       // client type, factory, context module and plugins
}

// PublisherData feeds the template of the publisher inject builds.
type PublisherData struct {
	Key     string          // config key naming the driver
	Default string          // driver when the key is unset
	Names   string          // driver names, comma separated
	Cases   []PublisherCase // sorted by name
	Imports []ImportData    // packages of the drivers and clients built anywhere
}

// PublisherCase is one driver of the publisher and the fx options wiring it.
type PublisherCase struct {
	Name    string
	Options string
}

// PublisherDriverData feeds the templates of a driver whose client builds
// only under a constraint: its file and the fallback without the constraint.
type PublisherDriverData struct {
	Name       string       // driver name
	Constraint string       // build constraint of its client
	Options    string       // fx options wiring the driver and its client
	Imports    []ImportData // packages of the driver and its client
}

// ClientPluginData is one plugin of a client.
type ClientPluginData struct {
	Name     string // plugin name, the last segment of its config root
	Alias    string // import alias of the plugin
	New      string // ...WithConfigPath constructor of the plugin
	Register string // method of the plugin passed to the factory
}

// ImplementationCase is one implementation of a switch and its module call.
type ImplementationCase struct {
	Name   string
	Module string
}

// ImportData is an aliased import of a generated file.
type ImportData struct {
	Alias string
	Path  string
}

// NewTemplate parses the generator templates ("module", "root",
// "implementation", "config", "client", "publisher", "publisher_driver" and
// "publisher_unsupported").
func NewTemplate() (*template.Template, error) {
	tmpl, err := template.New("inject").Parse(moduleTemplate)
	if err != nil {
		return nil, errors.Wrap(err, errors.Internalf("error parsing template"))
	}
	return tmpl, nil
}
