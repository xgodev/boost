package builder

import (
	"fmt"
	"go/types"
	"path"
	"strings"

	"github.com/xgodev/boost/inject"
)

// Import paths and aliases of the boost packages the decorators use.
const (
	middlewarePath     = "github.com/xgodev/boost/extra/middleware"
	middlewareFactPath = "github.com/xgodev/boost/factory/local/extra/middleware"
)

// decoratorData feeds the template of a decorator file
// (builder/middleware/<import path>/<func>_middleware.go).
type decoratorData struct {
	Header      string
	Package     string              // package name of the constructor
	Imports     []inject.ImportData // imports besides fx
	Func        string              // the constructor's function name: Decorate<Func>, decorate<Func>
	Constructor string              // pkg.Func, for the docs
	Interface   string              // pkg.Interface, for the docs
	Plugins     string              // plugin names, for the docs
	Decorate    string              // argument of fx.Decorate
	Type        string              // decorator type
	Next        string              // qualified interface type
	Context     string              // alias of the context package
	Fact        string              // alias of factory/local/extra/middleware
	// Gate is the call returning the chosen implementation when the
	// constructor is a member of an implementation set, "" otherwise; the
	// decorator applies only when it returns Member.
	Gate     string
	Member   string
	Wrappers []wrapperData
}

// wrapperData is one method of the decorated interface and the wrapper that
// runs it.
type wrapperData struct {
	Method  string
	Field   string   // decorator field holding the wrapper
	Name    string   // wrapper name
	Type    string   // wrapper type
	New     string   // wrapper constructor
	Plugins []string // plugin constructors, the outermost first
	Params  string   // parameters after ctx, with a leading ", "
	Results string   // results, as written in a signature
	Args    string   // arguments after ctx, with a leading ", "
	Return  string   // parameters of the identity return func
	Values  string   // what the identity return func returns
}

// decoratorFile returns the path, relative to OutputDir, of the decorator of
// the constructor m decorates, given the import path of the output dir.
func decoratorFile(base string, m Middleware) string {
	entry := m.Provider.Entry
	return path.Join(inject.GenPathUnder(base, path.Join(outputDir, middlewareDir), entry.Path), strings.ToLower(entry.Func.Name)+"_middleware.go")
}

// decoratorOf returns the import of the decorator package of m, given the
// import path of the output dir, and the call that applies it in Modules().
func decoratorOf(base string, m Middleware) (inject.ImportData, string) {
	entry := m.Provider.Entry
	p := path.Join(base, inject.GenPathUnder(base, path.Join(outputDir, middlewareDir), entry.Path))
	imp := inject.ImportData{Alias: inject.Alias("mw", entry.Package, p), Path: p}
	return imp, imp.Alias + ".Decorate" + entry.Func.Name + "()"
}

// decorator returns the data of the decorator file of m, given the import
// path of the output dir.
func decorator(base string, m Middleware) decoratorData {
	entry := m.Provider.Entry
	imps := newImports()
	ctxPkg := imps.add("context", contextPath)
	qualifier := func(pkg *types.Package) string {
		if pkg.Path() == contextPath {
			return ctxPkg
		}
		return imps.source(pkg.Name(), pkg.Path())
	}
	typeString := func(t types.Type) string { return types.TypeString(t, qualifier) }

	iface := m.Interface.Obj()
	d := decoratorData{
		Header:      generatedHeader,
		Package:     entry.Package,
		Func:        entry.Func.Name,
		Constructor: entry.Package + "." + entry.Func.Name,
		Interface:   iface.Pkg().Name() + "." + iface.Name(),
		Plugins:     strings.Join(m.Plugins, ", "),
		Decorate:    "decorate" + entry.Func.Name,
		Type:        strings.ToLower(entry.Func.Name[:1]) + entry.Func.Name[1:] + "Decorator",
		Next:        typeString(m.Interface),
		Context:     ctxPkg,
		Fact:        imps.add("fmiddleware", middlewareFactPath),
	}
	if name := m.Provider.Provides[m.Index].Name; name != "" {
		tag := goTag(`name:"` + name + `"`)
		d.Decorate = fmt.Sprintf("fx.Annotate(%s, fx.ParamTags(``, %s), fx.ResultTags(%s))", d.Decorate, tag, tag)
	}
	imps.add("contextfx", contextFxPath)
	if m.Implementation != nil {
		imp, call := m.Implementation.Chosen(base)
		d.Gate = imps.add(imp.Alias, imp.Path) + strings.TrimPrefix(call, imp.Alias)
		d.Member = entry.Package
	}

	mw := imps.add("middleware", middlewarePath)
	plugins := make([]string, len(m.Plugins))
	for i, name := range m.Plugins {
		plugin := middlewarePlugins[name]
		plugins[i] = imps.add(plugin.Alias, plugin.Path)
	}

	underlying := m.Interface.Underlying().(*types.Interface)
	for i := range underlying.NumMethods() {
		method := underlying.Method(i)
		sig := method.Signature()
		w := wrapperData{
			Method: method.Name(),
			Field:  strings.ToLower(method.Name()[:1]) + method.Name()[1:] + "Wrapper",
			Name:   wrapperName(m.Interface, method),
		}

		params := sig.Params()
		for j := 1; j < params.Len(); j++ {
			arg := fmt.Sprintf("a%d", j)
			tp := typeString(params.At(j).Type())
			if sig.Variadic() && j == params.Len()-1 {
				tp = "..." + typeString(params.At(j).Type().(*types.Slice).Elem())
				arg += "..."
			}
			w.Params += fmt.Sprintf(", a%d %s", j, tp)
			w.Args += ", " + arg
		}

		kind, _ := wrapperOf(sig)
		var constructor string
		switch kind {
		case errorWrapper:
			w.Type, w.New, constructor = "*"+mw+".ErrorWrapper", d.Fact+".NewErrorWrapperWithOptions", "NewError"
			w.Results, w.Return, w.Values = "error", "err error", "err"
		case anyErrorWrapper:
			r := typeString(sig.Results().At(0).Type())
			w.Type, w.New, constructor = "*"+mw+".AnyErrorWrapper["+r+"]", d.Fact+".NewAnyErrorWrapperWithOptions", "NewAnyError["+r+"]"
			w.Results, w.Return, w.Values = "("+r+", error)", "r "+r+", err error", "r, err"
		case anyWrapper:
			r := typeString(sig.Results().At(0).Type())
			w.Type, w.New, constructor = "*"+mw+".AnyWrapper["+r+"]", d.Fact+".NewAnyWrapperWithOptions", "NewAny["+r+"]"
			w.Results, w.Return, w.Values = r, "r "+r, "r"
		}
		for _, plugin := range plugins {
			w.Plugins = append(w.Plugins, plugin+"."+constructor)
		}
		d.Wrappers = append(d.Wrappers, w)
	}

	d.Imports = imps.sorted()
	return d
}
