package builder

import (
	"go/types"
	"strings"

	"github.com/xgodev/boost/annotation"
	"github.com/xgodev/boost/inject"
	"github.com/xgodev/boost/model/errors"
	"golang.org/x/tools/go/packages"
)

const middlewarePluginsRoot = "github.com/xgodev/boost/factory/local/extra/middleware/plugins"

// middlewarePlugins maps the plugin names of @Middleware to the package of
// their constructors (NewAnyError[R], NewAny[R] and NewError), imported under
// fixed aliases by the generated decorators.
var middlewarePlugins = map[string]inject.ImportData{
	"log":        {Alias: "logmw", Path: middlewarePluginsRoot + "/local/wrapper/log"},
	"hystrix":    {Alias: "hystrixmw", Path: middlewarePluginsRoot + "/contrib/afex/hystrix-go/v0"},
	"datadog":    {Alias: "datadogmw", Path: middlewarePluginsRoot + "/contrib/datadog/dd-trace-go/v1"},
	"prometheus": {Alias: "prometheusmw", Path: middlewarePluginsRoot + "/contrib/prometheus/client_golang/v1"},
}

// PackageLoader returns the type-checked packages at the import paths, keyed
// by path. NewSpec uses it to read the interfaces @Middleware decorates.
type PackageLoader func(paths ...string) (map[string]*types.Package, error)

// LoadPackages returns a PackageLoader that loads the packages with
// go/packages from dir, which must be inside the module that imports them.
func LoadPackages(dir string) PackageLoader {
	return func(paths ...string) (map[string]*types.Package, error) {
		mode := packages.NeedName | packages.NeedTypes | packages.NeedImports | packages.NeedDeps
		pkgs, err := packages.Load(&packages.Config{Mode: mode, Dir: dir}, paths...)
		if err != nil {
			return nil, errors.Wrap(err, errors.Internalf("error loading %s", strings.Join(paths, ", ")))
		}
		loaded := make(map[string]*types.Package, len(pkgs))
		for _, pkg := range pkgs {
			if len(pkg.Errors) > 0 {
				return nil, errors.NotValidf("package %s does not load: %v", pkg.PkgPath, pkg.Errors[0])
			}
			loaded[pkg.PkgPath] = pkg.Types
		}
		return loaded, nil
	}
}

// addMiddleware records a @Middleware whose interface is read later by
// loadMiddlewares, once for every annotated package.
func (b *specBuilder) addMiddleware(entry annotation.Entry, ann annotation.Annotation) error {
	where := declName(entry)
	if entry.IsStruct() {
		return errors.NotValidf("@Middleware on %s: only a @Provide function can be decorated", where)
	}
	key := inject.Component{Entry: entry}.Key()
	if _, repeated := b.middlewares[key]; repeated {
		return errors.NotValidf("@Middleware on %s is repeated; list every plugin in one @Middleware", where)
	}

	plugins, err := middlewarePluginNames(ann, where)
	if err != nil {
		return err
	}

	vertex := b.graph.Vertices[key]
	if vertex == nil || vertex.Value.Type != inject.AnnotationTypePROVIDE || len(vertex.Value.Provides) != 1 {
		return errors.NotValidf("@Middleware on %s needs exactly one @Provide(index=N) on the same function", where)
	}
	var index int
	for idx, provide := range vertex.Value.Provides {
		if provide.Group != "" {
			return errors.NotValidf("@Middleware on %s: @Provide(group=%s) is not supported, a decorator of a group receives the whole group", where, provide.Group)
		}
		index = idx
	}

	m := &Middleware{Provider: vertex.Value, Index: index, Plugins: plugins}
	if set, ok := inject.ImplementationOf(b.graph, vertex.Value); ok {
		m.Implementation = &set
	}
	b.middlewares[key] = m
	return nil
}

// middlewarePluginNames returns the plugin names of @Middleware(log, hystrix).
func middlewarePluginNames(ann annotation.Annotation, where string) ([]string, error) {
	valid := strings.Join(sortedNames(middlewarePlugins), ", ")
	var names []string
	seen := map[string]bool{}
	for _, name := range strings.Split(ann.RawValue(), ",") {
		name = strings.TrimSpace(name)
		switch {
		case name == "":
			continue
		case strings.Contains(name, "="):
			return nil, errors.NotValidf("@Middleware on %s takes plugin names, not %q; valid names: %s", where, name, valid)
		case seen[name]:
			return nil, errors.NotValidf("@Middleware on %s lists %s twice", where, name)
		}
		if _, ok := middlewarePlugins[name]; !ok {
			return nil, errors.NotValidf("@Middleware on %s: plugin %s is not one of %s", where, name, valid)
		}
		seen[name] = true
		names = append(names, name)
	}
	if len(names) == 0 {
		return nil, errors.NotValidf("@Middleware on %s needs at least one plugin: %s", where, valid)
	}
	return names, nil
}

// loadMiddlewares reads the interface each @Middleware decorates, with one
// call to load, and checks that the decorator can implement it.
func (b *specBuilder) loadMiddlewares(load PackageLoader) error {
	if len(b.middlewares) == 0 {
		return nil
	}
	if load == nil {
		return errors.NotValidf("@Middleware needs a PackageLoader to read the interfaces it decorates")
	}

	paths := map[string]bool{}
	for _, m := range b.middlewares {
		paths[m.Provider.Entry.Path] = true
	}
	pkgs, err := load(sortedNames(paths)...)
	if err != nil {
		return err
	}

	wrappers := map[string]*Middleware{} // wrapper name -> middleware that uses it
	for _, key := range sortedNames(b.middlewares) {
		m := b.middlewares[key]
		entry := m.Provider.Entry
		where := declName(entry)

		pkg := pkgs[entry.Path]
		if pkg == nil {
			return errors.NotFoundf("@Middleware on %s: package %s was not loaded", where, entry.Path)
		}
		fn, ok := pkg.Scope().Lookup(entry.Func.Name).(*types.Func)
		if !ok {
			return errors.NotFoundf("@Middleware on %s: package %s has no function %s", where, entry.Path, entry.Func.Name)
		}

		named, err := decoratedInterface(where, fn.Signature().Results().At(m.Index).Type())
		if err != nil {
			return err
		}
		iface := named.Underlying().(*types.Interface)
		for i := range iface.NumMethods() {
			name := wrapperName(named, iface.Method(i))
			if other, used := wrappers[name]; used && !sameImplementation(m, other) {
				return errors.AlreadyExistsf("@Middleware on %s: wrapper %s is already used by %s; decorate %s once",
					where, name, declName(other.Provider.Entry), types.TypeString(named, nil))
			}
			wrappers[name] = m
		}
		m.Interface = named
	}
	return nil
}

// sameImplementation reports whether a and b decorate members of the same
// implementation set: config wires one of them, so they may share wrappers.
func sameImplementation(a, b *Middleware) bool {
	if a.Implementation == nil || b.Implementation == nil {
		return false
	}
	x, y := a.Implementation, b.Implementation
	return x.Path == y.Path && x.Type == y.Type && x.Name == y.Name
}

// decoratedInterface checks that t is an interface a generated decorator, in
// another package, can implement and returns it.
func decoratedInterface(where string, t types.Type) (*types.Named, error) {
	named, ok := types.Unalias(t).(*types.Named)
	if !ok || !types.IsInterface(named) {
		return nil, errors.NotValidf("@Middleware on %s: %s is not a named interface type", where, types.TypeString(t, nil))
	}
	iface := named.Underlying().(*types.Interface)
	switch {
	case named.TypeParams().Len() > 0 || named.TypeArgs().Len() > 0:
		return nil, errors.NotValidf("@Middleware on %s: generic interface %s is not supported", where, types.TypeString(t, nil))
	case !named.Obj().Exported() || !iface.IsMethodSet():
		return nil, errors.NotValidf("@Middleware on %s: %s must be an exported interface of methods", where, types.TypeString(t, nil))
	case iface.NumMethods() == 0:
		return nil, errors.NotValidf("@Middleware on %s: %s has no methods to decorate", where, types.TypeString(t, nil))
	}

	for i := range iface.NumMethods() {
		if err := decoratedMethod(where, named, iface.Method(i)); err != nil {
			return nil, err
		}
	}
	return named, nil
}

// decoratedMethod checks that a decorator in another package can implement
// method of iface and run it in a wrapper.
func decoratedMethod(where string, iface *types.Named, method *types.Func) error {
	what := types.TypeString(iface, nil) + "." + method.Name()
	if !method.Exported() {
		return errors.NotValidf("@Middleware on %s: method %s is unexported, so a decorator in another package cannot implement it", where, what)
	}
	sig := method.Signature()
	if sig.Params().Len() == 0 || !isContext(sig.Params().At(0).Type()) {
		return errors.NotValidf("@Middleware on %s: method %s must take context.Context first", where, what)
	}
	if _, ok := wrapperOf(sig); !ok {
		return errors.NotValidf("@Middleware on %s: method %s must return error, (R, error) or R", where, what)
	}
	for _, tuple := range []*types.Tuple{sig.Params(), sig.Results()} {
		for v := range tuple.Variables() {
			if !exportedType(v.Type()) {
				return errors.NotValidf("@Middleware on %s: method %s uses unexported type %s", where, what, types.TypeString(v.Type(), nil))
			}
		}
	}
	return nil
}

// wrapperKind is the extra/middleware wrapper that runs a method.
type wrapperKind int

const (
	errorWrapper    wrapperKind = iota // func(...) error
	anyErrorWrapper                    // func(...) (R, error)
	anyWrapper                         // func(...) R
)

// wrapperOf tells which wrapper runs a method with signature sig.
func wrapperOf(sig *types.Signature) (wrapperKind, bool) {
	results := sig.Results()
	switch {
	case results.Len() == 1 && isError(results.At(0).Type()):
		return errorWrapper, true
	case results.Len() == 1:
		return anyWrapper, true
	case results.Len() == 2 && isError(results.At(1).Type()):
		return anyErrorWrapper, true
	}
	return 0, false
}

// wrapperName names the wrapper of a method of iface: the hystrix command,
// the span and metric name and the config path of every plugin.
func wrapperName(iface *types.Named, method *types.Func) string {
	return strings.ToLower(iface.Obj().Pkg().Name() + "." + iface.Obj().Name() + "." + method.Name())
}

func isContext(t types.Type) bool {
	named, ok := types.Unalias(t).(*types.Named)
	return ok && named.Obj().Pkg() != nil && named.Obj().Pkg().Path() == contextPath && named.Obj().Name() == "Context"
}

func isError(t types.Type) bool {
	return types.Identical(t, types.Universe.Lookup("error").Type())
}

// exportedType reports whether code in another package can write t.
func exportedType(t types.Type) bool {
	switch t := t.(type) {
	case *types.Alias:
		return exportedObj(t.Obj()) && exportedTypes(t.TypeArgs())
	case *types.Named:
		return exportedObj(t.Obj()) && exportedTypes(t.TypeArgs())
	case *types.Map:
		return exportedType(t.Key()) && exportedType(t.Elem())
	case interface{ Elem() types.Type }: // pointer, slice, array, chan
		return exportedType(t.Elem())
	case *types.Signature:
		return exportedTuple(t.Params()) && exportedTuple(t.Results())
	case *types.Struct:
		return exportedStruct(t)
	case *types.Interface:
		return exportedInterface(t)
	}
	return true
}

func exportedStruct(t *types.Struct) bool {
	for field := range t.Fields() {
		if !field.Exported() || !exportedType(field.Type()) {
			return false
		}
	}
	return true
}

func exportedInterface(t *types.Interface) bool {
	for i := range t.NumEmbeddeds() {
		if !exportedType(t.EmbeddedType(i)) {
			return false
		}
	}
	for i := range t.NumExplicitMethods() {
		method := t.ExplicitMethod(i)
		if !method.Exported() || !exportedType(method.Type()) {
			return false
		}
	}
	return true
}

func exportedObj(obj *types.TypeName) bool {
	return obj.Pkg() == nil || obj.Exported()
}

func exportedTypes(list *types.TypeList) bool {
	for t := range list.Types() {
		if !exportedType(t) {
			return false
		}
	}
	return true
}

func exportedTuple(tuple *types.Tuple) bool {
	for v := range tuple.Variables() {
		if !exportedType(v.Type()) {
			return false
		}
	}
	return true
}

// reached returns the keys of the components the inject modules of roots
// include: the roots and, transitively, the providers they depend on.
func reached(roots []inject.Component, vertices func(key string) []string) map[string]bool {
	seen := map[string]bool{}
	stack := make([]string, 0, len(roots))
	for _, root := range roots {
		stack = append(stack, root.Key())
	}
	for len(stack) > 0 {
		key := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if seen[key] {
			continue
		}
		seen[key] = true
		stack = append(stack, vertices(key)...)
	}
	return seen
}

// appRoots returns the inject components the doors of app include.
func appRoots(app App) []inject.Component {
	var roots []inject.Component
	receiver := func(h Handler) {
		if h.Receiver != nil {
			roots = append(roots, h.Receiver.Provider)
		}
	}
	for _, srv := range app.Servers {
		for _, route := range srv.Routes {
			receiver(route.Handler)
		}
		for _, svc := range srv.Services {
			roots = append(roots, svc.Provider)
		}
		roots = append(roots, srv.Plugins...)
	}
	if app.Function != nil {
		receiver(app.Function.Handler)
	}
	if app.Job != nil {
		receiver(app.Job.Handler)
	}
	return roots
}

// appMiddlewares returns the middlewares of the constructors in keys, sorted
// by constructor.
func appMiddlewares(middlewares map[string]*Middleware, keys map[string]bool) []Middleware {
	var list []Middleware
	for _, key := range sortedNames(keys) {
		if m, ok := middlewares[key]; ok {
			list = append(list, *m)
		}
	}
	return list
}
