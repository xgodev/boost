package builder

import (
	"fmt"
	"strconv"
	"strings"

	functionfx "github.com/xgodev/boost/fx/modules/bootstrap/function"
	"github.com/xgodev/boost/inject"
)

func functionDoor(base, app string, fn Function) (doorData, error) {
	d := door(app, functionFile, fmt.Sprintf("runs the @Function handler with the %s adapter.", fn.Adapter))
	imps := newImports()
	entry := fn.Handler.Entry

	result, err := imps.qualify(entry, entry.Func.Results[0].Type)
	if err != nil {
		return doorData{}, err
	}

	var params []string
	var modules []string
	include := func(provider inject.Component) {
		imp, mod := inject.ModuleOf(base, provider.Entry)
		imps.add(imp.Alias, imp.Path)
		modules = append(modules, mod)
	}
	call, err := handlerFunc(imps, fn.Handler, map[string]string{}, &params, include)
	if err != nil {
		return doorData{}, err
	}
	// the middlewares are built with the values their constructors provide
	for _, provider := range fn.Providers {
		include(provider)
	}

	handlerType := imps.add("bsfunction", functionPath) + ".Handler[" + result + "]"
	d.Options = append(d.Options, modules...)
	d.Options = append(d.Options,
		"fx.Provide(func("+strings.Join(params, ", ")+") "+handlerType+" {\nreturn "+call+"\n})",
	)
	d.Options = append(d.Options, functionMiddlewareOptions(imps, result, fn.Middlewares)...)
	d.Options = append(d.Options,
		imps.add("functionfx", functionFxPath)+".Module["+result+"](nil)",
	)
	if len(fn.Alternatives) == 0 {
		options, err := adapterOptions(base, imps, result, fn, fn.Adapter)
		if err != nil {
			return doorData{}, err
		}
		d.Options = append(d.Options, options...)
		d.Imports = imps.sorted()
		return d, nil
	}

	// a list of adapters: config chooses one, and only its module is wired
	names := fn.Adapters()
	sw := &adapterSwitch{
		Key:    configName(app) + "." + functionFile + ".adapter",
		Names:  strings.Join(names, ", "),
		Config: imps.add("config", wrapperConfigPath),
		Errors: imps.add("errors", errorsPath),
	}
	for _, name := range names {
		options, err := adapterOptions(base, imps, result, fn, name)
		if err != nil {
			return doorData{}, err
		}
		sw.Cases = append(sw.Cases, adapterCase{Name: name, Options: options})
	}
	d.Adapter = sw
	d.Doc = fmt.Sprintf("runs the @Function handler with the adapter named by %s (%s; default %s).", sw.Key, sw.Names, fn.Adapter)
	d.Init = []string{fmt.Sprintf("%s.Add(%sAdapterKey, %q, %q)", sw.Config, d.Var, fn.Adapter,
		"adapter of the @Function of app "+app+": "+sw.Names)}
	d.Options = append(d.Options, d.Var+"AdapterModule()")
	d.Imports = imps.sorted()
	return d, nil
}

// adapterOptions returns the fx modules that run the named adapter of fn and
// the plugins it is built with.
// An adapter reading with a client inject builds takes the app's client of
// that type, the default one or the one client=x names: the module of its
// inject constructor, or the one inject generates.
func adapterOptions(base string, imps *imports, result string, fn Function, name string) ([]string, error) {
	var options []string
	for _, adapter := range functionAdapters[name] {
		alias := imps.add(adapter.Alias, adapter.Path)
		if adapter.Client == "" {
			options = append(options, alias+".Module["+result+"]()")
			continue
		}

		var imp inject.ImportData
		var call string
		if provider, ok := fn.ClientProviders[adapter.Client]; ok {
			imp, call = inject.ModuleOf(base, provider.Entry)
		} else {
			var err error
			if imp, call, err = inject.NamedClientModule(base, adapter.Client, fn.Client); err != nil {
				return nil, err
			}
		}
		imps.add(imp.Alias, imp.Path)
		module := alias + ".ModuleWithoutClient[" + result + "]()"
		if fn.Client != "" {
			module = alias + ".ModuleWithNamedClient[" + result + "](" + strconv.Quote(fn.Client) + ")"
		}
		options = append(options, module, call)
	}
	// the adapter is built with the plugins its @Server lists
	if plugins := fn.Builtin[name]; len(plugins) > 0 {
		options = append(options, builtinOption(imps, name, plugins, ""))
	}
	return options, nil
}

// functionMiddlewareOptions returns, for the middlewares of a @Function, the
// fx module of each one and the provider of the list the function module runs,
// in the given order. A middleware disabled in config is left out.
func functionMiddlewareOptions(imps *imports, result string, names []string) []string {
	if len(names) == 0 {
		return nil
	}

	list := "[]" + imps.add("bsmiddleware", middlewarePath) + ".AnyErrorMiddleware[" + result + "]"
	var options, params, tags, appends []string
	for i, name := range names {
		mw := functionMiddlewares[name]
		options = append(options, imps.add(mw.Module.Alias, mw.Module.Path)+".Module["+result+"]()")
		param := fmt.Sprintf("m%d", i)
		params = append(params, param+" *"+imps.add(mw.Pkg.Alias, mw.Pkg.Path)+"."+mw.Type+"["+result+"]")
		tags = append(tags, goTag(`optional:"true"`))
		appends = append(appends, "if "+param+" != nil {\nm = append(m, "+param+")\n}\n")
	}

	return append(options, "fx.Provide(fx.Annotate(\n"+
		"func("+strings.Join(params, ", ")+") "+list+" {\n"+
		"// a middleware disabled in config is not provided\n"+
		"var m "+list+"\n"+
		strings.Join(appends, "")+
		"return m\n},\n"+
		"fx.ParamTags("+strings.Join(tags, ", ")+"),\n"+
		"fx.ResultTags("+goTag(`name:"`+functionfx.MiddlewaresName+`"`)+"),\n))")
}
