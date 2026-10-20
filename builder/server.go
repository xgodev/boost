package builder

import (
	"slices"
	"sort"
	"strings"

	"github.com/xgodev/boost/annotation"
	"github.com/xgodev/boost/inject"
	"github.com/xgodev/boost/model/errors"
)

// serverArgs are the arguments of
// @Server(type=rest, name=rest, app=api, plugins=[recover, logger]).
type serverArgs struct {
	Type    string `attr:"type"`
	Name    string `attr:"name"`
	App     string `attr:"app"`
	Plugins string `attr:"plugins"`
}

var serverKeys = []string{"type", "name", "app", "plugins"}

// serverType is a type of @Server: a server the builder generates, the
// package of the Plugin type it is built with and the plugins @Server lists.
// A new server type is a new entry of serverTypes.
type serverType struct {
	// Pkg is the package declaring the Plugin type.
	Pkg string
	// Alias prefixes the plugin packages imported by the generated door.
	Alias string
	// Plugins maps each plugin name to the package of its Register
	// function. The name is the one the plugin's enabled key has in config,
	// <plugins root>.<name>.enabled, so config turns off what @Server lists.
	Plugins map[string]string
	// Defaults are the plugins of a server of the type without @Server, in
	// order.
	Defaults []string
	// Adapter is the @Function adapter the type is, "" for a server of its
	// own: its plugins go to the adapter of the app's function and name=
	// does not apply.
	Adapter string
}

const (
	echoPlugins = boostEchoPath + "/plugins"
	grpcPlugins = grpcServerPath + "/plugins"
	mcpPlugins  = mcpServerPath + "/plugins"
	ceHTTPPath  = functionPath + "/adapter/contrib/cloudevents/sdk-go/v2/core/http"
)

// serverTypes are the types of @Server, by the name type= takes.
var serverTypes = map[string]serverType{
	string(KindREST): {
		Pkg:   boostEchoPath,
		Alias: "echo",
		Plugins: map[string]string{
			"bodydump":       echoPlugins + "/native/bodydump",
			"bodyLimit":      echoPlugins + "/native/bodylimit",
			"bytedanceSonic": echoPlugins + "/contrib/bytedance/sonic/v1",
			"cors":           echoPlugins + "/native/cors",
			"datadog":        echoPlugins + "/contrib/datadog/dd-trace-go/v1",
			"errorHandler":   echoPlugins + "/extra/error_handler",
			"goccyGoJson":    echoPlugins + "/contrib/goccy/go-json/v0",
			"gzip":           echoPlugins + "/native/gzip",
			"health":         echoPlugins + "/local/extra/health",
			"logger":         echoPlugins + "/local/wrapper/log",
			"multiServer":    echoPlugins + "/local/extra/multi-server",
			"otel":           echoPlugins + "/contrib/go.opentelemetry.io/contrib/v0",
			"pprof":          echoPlugins + "/contrib/hiko1129/echo-pprof/v1",
			"prometheus":     echoPlugins + "/contrib/prometheus/client_golang/v1",
			"recover":        echoPlugins + "/native/recover",
			"requestId":      echoPlugins + "/native/requestid",
			"semaphore":      echoPlugins + "/extra/semaphore",
			"status":         echoPlugins + "/local/model/restresponse",
			"swagger":        echoPlugins + "/contrib/swaggo/echo-swagger/v1",
		},
		Defaults: []string{"recover", "requestId", "logger", "errorHandler", "health"},
	},
	string(KindGRPC): {
		Pkg:   grpcServerPath,
		Alias: "grpc",
		Plugins: map[string]string{
			"compressor": grpcPlugins + "/native/compressor",
			"datadog":    grpcPlugins + "/contrib/datadog/dd-trace-go/v1",
			"log":        grpcPlugins + "/local/wrapper/log",
			"prometheus": grpcPlugins + "/contrib/prometheus/client_golang/v1",
		},
		Defaults: []string{"log"},
	},
	string(KindMCP): {
		Pkg:   mcpServerPath,
		Alias: "mcp",
		Plugins: map[string]string{
			"log":  mcpPlugins + "/local/wrapper/log",
			"otel": mcpPlugins + "/contrib/go.opentelemetry.io/otel/v1",
		},
		Defaults: []string{"log"},
	},
	"cloudevents": {
		Pkg:   ceHTTPPath,
		Alias: "ce",
		Plugins: map[string]string{
			"otel": ceHTTPPath + "/plugins/contrib/go.opentelemetry.io/contrib/v0",
		},
		Adapter: "cloudevents",
	},
}

// serverDecl is a @Server: the plugins of one server, or of one adapter of
// the function of an app.
type serverDecl struct {
	Type    string
	Plugins []string
	Where   string
	used    bool
}

// serverKey keys the @Server of server name of app, or of the adapter of
// its function.
func serverKey(app, name string, tp serverType) string {
	if tp.Adapter != "" {
		return app + " function " + tp.Adapter
	}
	return app + " server " + name
}

func (b *specBuilder) addServer(entry annotation.Entry, ann annotation.Annotation) error {
	where := declName(entry)
	if !entry.IsPackage() {
		return errors.NotValidf("@Server on %s: write it in the package doc of the server's package (its doc.go), above the package clause", where)
	}
	var args serverArgs
	if err := decode(ann, where, serverKeys, &args); err != nil {
		return err
	}

	types := strings.Join(sortedNames(serverTypes), ", ")
	if args.Type == "" {
		return errors.NotValidf("@Server on %s: type is required (one of %s)", where, types)
	}
	tp, ok := serverTypes[args.Type]
	if !ok {
		return errors.NotValidf("@Server on %s: type %s is not one of %s", where, args.Type, types)
	}
	if tp.Adapter != "" && args.Name != "" {
		return errors.NotValidf("@Server on %s: %s is the adapter of the app's @Function, so name= does not apply", where, args.Type)
	}

	app, name := orDefault(args.App, defaultApp), args.Name
	if tp.Adapter == "" {
		name = orDefault(name, defaultServers[Kind(args.Type)].name)
	}
	if err := validateNames(ann.Name, where, app, args.Name); err != nil {
		return err
	}

	plugins := tp.Defaults
	if _, listed := ann.Map["plugins"]; listed {
		var err error
		if plugins, err = parseList(args.Plugins, "@Server on "+where, "plugins"); err != nil {
			return err
		}
	}
	for _, plugin := range plugins {
		if _, ok := tp.Plugins[plugin]; !ok {
			return errors.NotValidf("@Server on %s: plugin %s of %s is not one of %s", where, plugin, args.Type, strings.Join(sortedNames(tp.Plugins), ", "))
		}
	}

	key := serverKey(app, name, tp)
	if other, dup := b.servers[key]; dup {
		what := "server " + name
		if tp.Adapter != "" {
			what = "the " + tp.Adapter + " adapter"
		}
		return errors.AlreadyExistsf("@Server on %s and on %s both set the plugins of %s of app %s; keep one", other.Where, where, what, app)
	}
	b.servers[key] = &serverDecl{Type: args.Type, Plugins: plugins, Where: where}
	return nil
}

// serverBuiltin returns the plugins of srv from its @Server or, without
// one, the default set of its type, unless an @Provide []Plugin gives them.
func (b *specBuilder) serverBuiltin(srv Server, provided []inject.Component) ([]string, error) {
	tp := serverTypes[string(srv.Kind)]
	decl, ok := b.servers[serverKey(srv.App, srv.Name, tp)]
	if !ok {
		if len(provided) > 0 {
			return nil, nil
		}
		return tp.Defaults, nil
	}

	decl.used = true
	switch {
	case decl.Type != string(srv.Kind):
		return nil, errors.NotValidf("@Server(type=%s) on %s: server %s of app %s serves %s", decl.Type, decl.Where, srv.Name, srv.App, srv.Kind)
	case len(provided) > 0:
		return nil, errors.AlreadyExistsf("@Server on %s lists the plugins of server %s of app %s and @Provide gives them too (%s); keep @Server",
			decl.Where, srv.Name, srv.App, declNames(provided))
	}
	return decl.Plugins, nil
}

// functionBuiltin returns the plugins of each adapter of fn, from its
// @Server or the default set of the adapter's type, by adapter.
func (b *specBuilder) functionBuiltin(app string, fn *Function) map[string][]string {
	var builtin map[string][]string
	for _, name := range sortedNames(serverTypes) {
		tp := serverTypes[name]
		if tp.Adapter == "" || !slices.Contains(fn.Adapters(), tp.Adapter) {
			continue
		}
		plugins := tp.Defaults
		if decl, ok := b.servers[serverKey(app, "", tp)]; ok {
			decl.used = true
			plugins = decl.Plugins
		}
		if len(plugins) == 0 {
			continue
		}
		if builtin == nil {
			builtin = map[string][]string{}
		}
		builtin[tp.Adapter] = plugins
	}
	return builtin
}

// unusedServer fails for the first @Server no server or adapter took: its
// app or name is wrong.
func (b *specBuilder) unusedServer() error {
	keys := make([]string, 0, len(b.servers))
	for key := range b.servers {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		decl := b.servers[key]
		if decl.used {
			continue
		}
		app, rest, _ := strings.Cut(key, " ")
		if tp := serverTypes[decl.Type]; tp.Adapter != "" {
			return errors.NotFoundf("@Server(type=%s) on %s: app %s has no @Function with the %s adapter", decl.Type, decl.Where, app, tp.Adapter)
		}
		return errors.NotFoundf("@Server(type=%s) on %s: app %s has no server %s; set name= and app= to a server its doors declare",
			decl.Type, decl.Where, app, strings.TrimPrefix(rest, "server "))
	}
	return nil
}

// builtinPlugins returns the type and the value of the []Plugin of the
// server type listing plugins, importing their packages:
// []pkg.Plugin and []pkg.Plugin{a.Register, ...}.
func builtinPlugins(imps *imports, typeName string, plugins []string) (string, string) {
	tp := serverTypes[typeName]
	// the door already imports the server's package under its own alias
	list := "[]" + imps.add(tp.Alias+"plugin", tp.Pkg) + ".Plugin"
	var registers []string
	for _, plugin := range plugins {
		registers = append(registers, imps.add(tp.Alias+strings.ToLower(plugin), tp.Plugins[plugin])+".Register")
	}
	return list, list + "{\n" + strings.Join(registers, ",\n") + ",\n}"
}

// builtinOption returns the fx option providing the []Plugin of the server
// type listing plugins, under name or unnamed when name is "".
func builtinOption(imps *imports, typeName string, plugins []string, name string) string {
	list, value := builtinPlugins(imps, typeName, plugins)
	provider := "func() " + list + " {\nreturn " + value + "\n}"
	if name == "" {
		return "fx.Provide(" + provider + ")"
	}
	return "fx.Provide(\nfx.Annotate(\n" + provider + ",\nfx.ResultTags(" + goTag(`name:"`+name+`"`) + "),\n),\n)"
}
