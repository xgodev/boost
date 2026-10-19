package builder

import (
	"go/token"
	"path"
	"sort"
	"strings"

	"github.com/xgodev/boost/annotation"
	"github.com/xgodev/boost/extra/graph"
	"github.com/xgodev/boost/inject"
	"github.com/xgodev/boost/model/errors"
)

// NewSpec builds the doors declared by the @REST, @GRPC, @MCP and @Function
// annotations of entries, and the @Middleware of the constructors each app
// builds. g is the inject graph of the same entries: it provides the
// receivers of handler methods and the values of @GRPC constructors. load
// reads the interfaces @Middleware decorates; it is only called, once, when
// there is a @Middleware.
//
// It returns NotValid for malformed annotations and for handlers that are
// unexported, in package main or whose signature does not match their door,
// NotFound when the receiver of a handler method or a value a @Function
// middleware needs has no provider, and
// AlreadyExists for a duplicate route, service or tool, for a second @Function in
// an app, for two @Middleware that decorate the same interface and for a rest,
// grpc or mcp server given both a named and an unnamed []Plugin.
func NewSpec(entries []annotation.Entry, g *graph.Graph[inject.Component], load PackageLoader) (*Spec, error) {
	b := &specBuilder{graph: g, apps: map[string]*appState{}, middlewares: map[string]*Middleware{}}

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
			case "MCP":
				err = b.addTool(entry, ann)
			case "Function":
				err = b.addFunction(entry, ann)
			case "Middleware":
				err = b.addMiddleware(entry, ann)
			}
			if err != nil {
				return nil, err
			}
		}
	}

	if err := b.loadMiddlewares(load); err != nil {
		return nil, err
	}
	return b.spec()
}

type specBuilder struct {
	graph       *graph.Graph[inject.Component]
	apps        map[string]*appState
	middlewares map[string]*Middleware // by constructor key
}

type appState struct {
	servers  map[string]*Server
	function *Function
	served   map[string]bool // routes, services and tools already declared, per server
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

// serve records a route, service or tool of a server, failing when it is declared twice.
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
	if err := b.singleConstructor(vertex.Value, ann.Name, where, "its constructor"); err != nil {
		return err
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

func (b *specBuilder) addTool(entry annotation.Entry, ann annotation.Annotation) error {
	where := declName(entry)
	var args mcpArgs
	if err := decode(ann, where, mcpKeys, &args); err != nil {
		return err
	}

	name := orDefault(args.Name, snakeCase(entry.Func.Name))
	if !toolName.MatchString(name) {
		return errors.NotValidf("@MCP on %s: tool name %q must match %s", where, name, toolName)
	}
	app, server := orDefault(args.App, defaultApp), orDefault(args.Server, defaultMCPServer)
	if err := validateNames(ann.Name, where, app, server); err != nil {
		return err
	}

	handler, err := b.handler(entry, ann.Name)
	if err != nil {
		return err
	}
	params, results := entry.Func.Parameters, entry.Func.Results
	if len(params) != 3 || !isType(entry, params[0].Type, typeRef{Path: "context", Name: "Context"}) ||
		!isPointer(entry, params[1].Type, typeRef{Path: mcpPath, Name: "CallToolRequest"}) ||
		len(results) != 3 || !isPointer(entry, results[0].Type, typeRef{Path: mcpPath, Name: "CallToolResult"}) ||
		results[2].Type != "error" {
		return errors.NotValidf("@MCP on %s: the handler must be func(context.Context, *mcp.CallToolRequest, In) (*mcp.CallToolResult, Out, error)", where)
	}

	srv, err := b.server(app, server, KindMCP, ann.Name, where)
	if err != nil {
		return err
	}
	if err := b.serve(app, server, "tool "+name, ann.Name, where); err != nil {
		return err
	}
	srv.Tools = append(srv.Tools, Tool{Name: name, Description: toolDescription(entry.Header), Handler: handler})
	return nil
}

func (b *specBuilder) addFunction(entry annotation.Entry, ann annotation.Annotation) error {
	where := declName(entry)
	var args functionArgs
	if err := decode(ann, where, functionKeys, &args); err != nil {
		return err
	}
	adapters, err := functionAdapterNames(args.Adapter, where)
	if err != nil {
		return err
	}

	middlewares, err := parseList(args.Middlewares, "@Function on "+where, "middlewares")
	if err != nil {
		return err
	}
	for _, name := range middlewares {
		if _, ok := functionMiddlewares[name]; !ok {
			return errors.NotValidf("@Function on %s: middleware %s is not one of %s", where, name, strings.Join(sortedNames(functionMiddlewares), ", "))
		}
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

	providers, err := middlewareProviders(b.graph, middlewares, where)
	if err != nil {
		return err
	}

	a := b.app(app)
	if a.function != nil {
		return errors.AlreadyExistsf("@Function on %s: app %s already runs %s and an app runs one @Function; set app= to another app",
			where, app, declName(a.function.Handler.Entry))
	}
	a.function = &Function{Adapter: adapters[0], Handler: handler, Middlewares: middlewares, Providers: providers,
		ClientProviders: adapterClientProviders(b.graph, adapters)}
	if len(adapters) > 1 {
		a.function.Alternatives = adapters[1:]
	}
	return nil
}

// adapterClientProviders returns the inject constructors of the default
// clients the adapters read with, by client type.
func adapterClientProviders(g *graph.Graph[inject.Component], adapters []string) map[string]inject.Component {
	providers := map[string]inject.Component{}
	for _, name := range adapters {
		for _, adapter := range functionAdapters[name] {
			if adapter.Client == "" {
				continue
			}
			if comp, ok := inject.DefaultClientProvider(g, adapter.Client); ok {
				providers[adapter.Client] = comp
			}
		}
	}
	return providers
}

// functionAdapterNames reads the adapter of @Function, one name (nats) or a
// list ([pubsub, nats]), and checks that every name is a known adapter.
func functionAdapterNames(raw, where string) ([]string, error) {
	names := []string{raw}
	if strings.HasPrefix(strings.TrimSpace(raw), "[") {
		var err error
		if names, err = parseList(raw, "@Function on "+where, "adapter"); err != nil {
			return nil, err
		}
	}

	valid := strings.Join(sortedNames(functionAdapters), ", ")
	if len(names) == 0 || names[0] == "" {
		return nil, errors.NotValidf("@Function on %s: adapter is required (one of %s)", where, valid)
	}
	for _, name := range names {
		if _, ok := functionAdapters[name]; !ok {
			return nil, errors.NotValidf("@Function on %s: adapter %s is not one of %s", where, name, valid)
		}
	}
	return names, nil
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
		if comp, _, ok := inject.ProviderFor(b.graph, entry, tp); ok {
			if err := b.singleConstructor(comp, ann, declName(entry), "its receiver "+tp); err != nil {
				return Handler{}, err
			}
			return Handler{Entry: entry, Receiver: &Receiver{Type: tp, Provider: comp}}, nil
		}
	}
	return Handler{}, errors.NotFoundf("@%s on %s: no @Provide returns *%s or %s, so the generated door cannot get the receiver",
		ann, declName(entry), entry.Struct, entry.Struct)
}

// singleConstructor fails when comp, the constructor a door is built from,
// is a member of an implementation set: the door would include that member
// whatever config chooses.
func (b *specBuilder) singleConstructor(comp inject.Component, ann, where, what string) error {
	set, ok := inject.ImplementationOf(b.graph, comp)
	if !ok {
		return nil
	}
	return errors.NotValidf("@%s on %s: %s comes from %s, one of the implementations of %s.%s (%s) chosen by config; a door is built from one constructor, so serve it from a type that receives the chosen %s",
		ann, where, what, comp.Key(), set.Path, set.Type, strings.Join(set.Names(), ", "), set.Type)
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
func (b *specBuilder) spec() (*Spec, error) {
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
			sort.SliceStable(srv.Tools, func(i, j int) bool { return srv.Tools[i].Name < srv.Tools[j].Name })
			plugins, pluginsName, err := serverPlugins(b.graph, srv)
			if err != nil {
				return nil, err
			}
			srv.Plugins, srv.PluginsName = plugins, pluginsName
			app.Servers = append(app.Servers, srv)
		}
		app.Middlewares = appMiddlewares(b.middlewares, reached(appRoots(app), b.dependencies))
		spec.Apps = append(spec.Apps, app)
	}
	return spec, nil
}

// dependencies returns the keys of the providers the component key depends on.
func (b *specBuilder) dependencies(key string) []string {
	vertex := b.graph.Vertices[key]
	if vertex == nil {
		return nil
	}
	var keys []string
	for _, dep := range vertex.Incoming() {
		keys = append(keys, dep.Key)
	}
	return keys
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

// middlewareProviders returns the inject constructors providing the unnamed
// values the @Function middlewares need. A value no constructor provides is
// NotFound: fx would skip the optional middleware without a word.
func middlewareProviders(g *graph.Graph[inject.Component], names []string, where string) ([]inject.Component, error) {
	var providers []inject.Component
	seen := map[string]bool{}
	for _, name := range names {
		for _, need := range functionMiddlewares[name].Needs {
			found := false
			for _, key := range g.SortedKeys() {
				comp := g.Vertices[key].Value
				for idx, ann := range comp.Provides {
					if ann.Group != "" || ann.Name != "" || !isType(comp.Entry, comp.Result(idx), need) {
						continue
					}
					found = true
					if !seen[key] {
						seen[key] = true
						providers = append(providers, comp)
					}
					break
				}
			}
			if !found {
				return nil, errors.NotFoundf("@Function on %s: middleware %s needs a %s.%s; annotate the constructor that returns it with @Provide",
					where, name, path.Base(need.Path), need.Name)
			}
		}
	}
	return providers, nil
}
