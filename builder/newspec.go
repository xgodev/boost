package builder

import (
	"go/ast"
	"go/parser"
	"go/token"
	"sort"
	"strings"

	"github.com/xgodev/boost/annotation"
	"github.com/xgodev/boost/extra/graph"
	"github.com/xgodev/boost/inject"
	"github.com/xgodev/boost/model/errors"
)

// NewSpec builds the doors declared by the @REST, @GRPC and @Function
// annotations of entries. g is the inject graph of the same entries: it
// provides the receivers of handler methods and the values of @GRPC
// constructors.
//
// It returns NotValid for malformed annotations and for handlers that are
// unexported, in package main or whose signature does not match their door,
// NotFound when the receiver of a handler method has no provider, and
// AlreadyExists for a duplicate route or service and for a second @Function
// in an app.
func NewSpec(entries []annotation.Entry, g *graph.Graph[inject.Component]) (*Spec, error) {
	b := &specBuilder{graph: g, apps: map[string]*appState{}}

	seen := map[string]bool{}
	for _, entry := range entries {
		// the collector may load a package twice
		if where := declName(entry); seen[where] {
			continue
		} else {
			seen[where] = true
		}

		for _, ann := range entry.Annotations {
			var err error
			switch ann.Name {
			case "REST":
				err = b.addRoute(entry, ann)
			case "GRPC":
				err = b.addService(entry, ann)
			case "Function":
				err = b.addFunction(entry, ann)
			}
			if err != nil {
				return nil, err
			}
		}
	}

	return b.spec(), nil
}

type specBuilder struct {
	graph *graph.Graph[inject.Component]
	apps  map[string]*appState
}

type appState struct {
	servers  map[string]*Server
	function *Function
	served   map[string]bool // routes and services already declared, per server
}

func (b *specBuilder) app(name string) *appState {
	a, ok := b.apps[name]
	if !ok {
		a = &appState{servers: map[string]*Server{}, served: map[string]bool{}}
		b.apps[name] = a
	}
	return a
}

// server returns the server of the app, failing when it already serves
// another kind of door.
func (b *specBuilder) server(app, name string, kind Kind, ann, where string) (*Server, error) {
	a := b.app(app)
	srv, ok := a.servers[name]
	if !ok {
		srv = &Server{App: app, Name: name, Kind: kind}
		a.servers[name] = srv
	}
	if srv.Kind != kind {
		return nil, errors.NotValidf("@%s on %s: server %s of app %s already serves %s; give the %s door another server name", ann, where, name, app, srv.Kind, kind)
	}
	return srv, nil
}

// serve records a route or service of a server, failing when it is declared twice.
func (b *specBuilder) serve(app, server, key, ann, where string) error {
	a := b.app(app)
	full := server + " " + key
	if a.served[full] {
		return errors.AlreadyExistsf("@%s on %s: %s is already served by server %s of app %s", ann, where, key, server, app)
	}
	a.served[full] = true
	return nil
}

func (b *specBuilder) addRoute(entry annotation.Entry, ann annotation.Annotation) error {
	where := declName(entry)
	var args restArgs
	if err := decode(ann, where, restKeys, &args); err != nil {
		return err
	}

	method := strings.ToUpper(args.Method)
	switch {
	case args.Method == "":
		return errors.NotValidf("@REST on %s: method is required", where)
	case !restMethods[method]:
		return errors.NotValidf("@REST on %s: method %s is not one of %s", where, args.Method, strings.Join(sortedNames(restMethods), ", "))
	case args.Path == "":
		return errors.NotValidf("@REST on %s: path is required", where)
	}

	app, server := orDefault(args.App, defaultApp), orDefault(args.Server, defaultRESTServer)
	if err := validateNames(ann.Name, where, app, server); err != nil {
		return err
	}

	handler, err := b.handler(entry, ann.Name)
	if err != nil {
		return err
	}
	params, results := entry.Func.Parameters, entry.Func.Results
	if len(params) != 1 || !isType(entry, params[0].Type, typeRef{Path: echoPath, Name: "Context"}) ||
		len(results) != 1 || results[0].Type != "error" {
		return errors.NotValidf("@REST on %s: the handler must be func(echo.Context) error", where)
	}

	srv, err := b.server(app, server, KindREST, ann.Name, where)
	if err != nil {
		return err
	}
	if err := b.serve(app, server, method+" "+args.Path, ann.Name, where); err != nil {
		return err
	}
	srv.Routes = append(srv.Routes, Route{Method: method, Path: args.Path, Handler: handler})
	return nil
}

func (b *specBuilder) addService(entry annotation.Entry, ann annotation.Annotation) error {
	where := declName(entry)
	var args grpcArgs
	if err := decode(ann, where, grpcKeys, &args); err != nil {
		return err
	}
	if args.Register == "" {
		return errors.NotValidf("@GRPC on %s: register is required", where)
	}

	app, server := orDefault(args.App, defaultApp), orDefault(args.Server, defaultGRPCServer)
	if err := validateNames(ann.Name, where, app, server); err != nil {
		return err
	}
	if err := checkDecl(entry, ann.Name); err != nil {
		return err
	}

	vertex := b.graph.Vertices[inject.Component{Entry: entry}.Key()]
	if vertex == nil || vertex.Value.Type != inject.AnnotationTypePROVIDE || len(vertex.Value.Provides) != 1 {
		return errors.NotValidf("@GRPC on %s needs exactly one @Provide(index=N) on the same function", where)
	}
	var index int
	for idx, provide := range vertex.Value.Provides {
		if provide.IDType() != inject.AnnotationIDTypeDEFAULT {
			return errors.NotValidf("@GRPC on %s needs @Provide(index=%d) without name or group", where, idx)
		}
		index = idx
	}

	ref, ok := resolve(entry, args.Register)
	if !ok || !token.IsExported(ref.Name) {
		return errors.NotValidf("@GRPC on %s: register %q does not resolve; write pkg.RegisterXServer with pkg imported by the file, or an exported function of the package", where, args.Register)
	}

	srv, err := b.server(app, server, KindGRPC, ann.Name, where)
	if err != nil {
		return err
	}
	if err := b.serve(app, server, ref.Path+"."+ref.Name, ann.Name, where); err != nil {
		return err
	}
	srv.Services = append(srv.Services, Service{Register: args.Register, Provider: vertex.Value, Index: index})
	return nil
}

func (b *specBuilder) addFunction(entry annotation.Entry, ann annotation.Annotation) error {
	where := declName(entry)
	var args functionArgs
	if err := decode(ann, where, functionKeys, &args); err != nil {
		return err
	}
	adapters := strings.Join(sortedNames(functionAdapters), ", ")
	if args.Adapter == "" {
		return errors.NotValidf("@Function on %s: adapter is required (one of %s)", where, adapters)
	}
	if _, ok := functionAdapters[args.Adapter]; !ok {
		return errors.NotValidf("@Function on %s: adapter %s is not one of %s", where, args.Adapter, adapters)
	}

	app := orDefault(args.App, defaultApp)
	if err := validateNames(ann.Name, where, app, ""); err != nil {
		return err
	}

	handler, err := b.handler(entry, ann.Name)
	if err != nil {
		return err
	}
	params, results := entry.Func.Parameters, entry.Func.Results
	if len(params) != 2 || !isType(entry, params[0].Type, typeRef{Path: "context", Name: "Context"}) ||
		!isType(entry, params[1].Type, typeRef{Path: cloudeventsPath, Name: "Event"}, typeRef{Path: eventPath, Name: "Event"}) ||
		len(results) != 2 || results[1].Type != "error" {
		return errors.NotValidf("@Function on %s: the handler must be func(context.Context, cloudevents.Event) (T, error)", where)
	}

	a := b.app(app)
	if a.function != nil {
		return errors.AlreadyExistsf("@Function on %s: app %s already runs %s and an app runs one @Function; set app= to another app",
			where, app, declName(a.function.Handler.Entry))
	}
	a.function = &Function{Adapter: args.Adapter, Handler: handler}
	return nil
}

// handler checks the declaration serving a door and finds the provider of
// its receiver when it is a method.
func (b *specBuilder) handler(entry annotation.Entry, ann string) (Handler, error) {
	if err := checkDecl(entry, ann); err != nil {
		return Handler{}, err
	}
	if !entry.IsMethod() {
		return Handler{Entry: entry}, nil
	}

	for _, tp := range []string{"*" + entry.Struct, entry.Struct} {
		if comp, _, ok := inject.Provider(b.graph, entry.Package, tp); ok {
			return Handler{Entry: entry, Receiver: &Receiver{Type: tp, Provider: comp}}, nil
		}
	}
	return Handler{}, errors.NotFoundf("@%s on %s: no @Provide returns *%s or %s, so the generated door cannot get the receiver",
		ann, declName(entry), entry.Struct, entry.Struct)
}

// checkDecl fails for declarations the generated code cannot call.
func checkDecl(entry annotation.Entry, ann string) error {
	where := declName(entry)
	switch {
	case entry.IsStruct():
		return errors.NotValidf("@%s on %s: only a function or method can serve a door", ann, where)
	case entry.Package == "main":
		return errors.NotValidf("@%s on %s: package main cannot be imported by the generated code; move it to another package", ann, where)
	case !token.IsExported(entry.Func.Name) || (entry.IsMethod() && !token.IsExported(entry.Struct)):
		return errors.NotValidf("@%s on %s: the handler must be exported to be called by the generated code", ann, where)
	}
	return nil
}

// spec sorts what was collected into the Spec.
func (b *specBuilder) spec() *Spec {
	spec := &Spec{}
	for _, name := range sortedNames(b.apps) {
		a := b.apps[name]
		app := App{Name: name, Function: a.function}
		for _, srvName := range sortedNames(a.servers) {
			srv := *a.servers[srvName]
			sort.SliceStable(srv.Routes, func(i, j int) bool {
				if srv.Routes[i].Path != srv.Routes[j].Path {
					return srv.Routes[i].Path < srv.Routes[j].Path
				}
				return srv.Routes[i].Method < srv.Routes[j].Method
			})
			sort.SliceStable(srv.Services, func(i, j int) bool {
				return srv.Services[i].Provider.Key()+srv.Services[i].Register < srv.Services[j].Provider.Key()+srv.Services[j].Register
			})
			srv.Plugins = serverPlugins(b.graph, srv)
			app.Servers = append(app.Servers, srv)
		}
		spec.Apps = append(spec.Apps, app)
	}
	return spec
}

// declName names a declaration in messages: path.Func, path.Type.Method or path.Type.
func declName(entry annotation.Entry) string {
	name := entry.Func.Name
	switch {
	case entry.Struct != "" && name != "":
		name = entry.Struct + "." + name
	case name == "":
		name = entry.Struct
	}
	return entry.Path + "." + name
}

// pluginPaths are the packages whose Plugin type a server of each kind is
// built with.
var pluginPaths = map[Kind]string{KindREST: boostEchoPath, KindGRPC: grpcServerPath}

// serverPlugins returns the inject constructors providing the []Plugin of
// srv: a default value for the rest and grpc servers, the value named after
// the server for the others.
func serverPlugins(g *graph.Graph[inject.Component], srv Server) []inject.Component {
	name := srv.Name
	if isDefault(srv) {
		name = ""
	}

	var plugins []inject.Component
	for _, key := range g.SortedKeys() {
		comp := g.Vertices[key].Value
		for idx, ann := range comp.Provides {
			if ann.Group == "" && ann.Name == name && isPluginSlice(comp.Entry, comp.Entry.Func.Results[idx].Type, pluginPaths[srv.Kind]) {
				plugins = append(plugins, comp)
				break
			}
		}
	}
	return plugins
}

// isPluginSlice reports whether tp, written in the file of entry, is
// []pkg.Plugin with pkg imported from pkgPath.
func isPluginSlice(entry annotation.Entry, tp, pkgPath string) bool {
	expr, err := parser.ParseExpr(tp)
	if err != nil {
		return false
	}
	slice, ok := expr.(*ast.ArrayType)
	if !ok || slice.Len != nil {
		return false
	}
	sel, ok := slice.Elt.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Plugin" {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	return ok && entry.Imports[pkg.Name] == pkgPath
}
