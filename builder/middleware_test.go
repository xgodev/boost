package builder

import (
	"context"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost/annotation"
	"github.com/xgodev/boost/inject"
	"github.com/xgodev/boost/model/errors"
)

// testLoader type-checks the sources, one file per import path, and loads
// them; other paths come from the standard library.
func testLoader(t *testing.T, sources map[string]string) PackageLoader {
	t.Helper()
	fset := token.NewFileSet()
	std := importer.Default()
	checked := map[string]*types.Package{}

	var imp types.ImporterFrom
	var check func(path string) (*types.Package, error)
	check = func(path string) (*types.Package, error) {
		if pkg, ok := checked[path]; ok {
			return pkg, nil
		}
		src, ok := sources[path]
		if !ok {
			return std.Import(path)
		}
		file, err := parser.ParseFile(fset, path+".go", src, 0)
		if err != nil {
			return nil, err
		}
		pkg, err := (&types.Config{Importer: imp}).Check(path, fset, []*ast.File{file}, nil)
		if err != nil {
			return nil, err
		}
		checked[path] = pkg
		return pkg, nil
	}
	imp = importerFunc(check)

	return func(paths ...string) (map[string]*types.Package, error) {
		loaded := map[string]*types.Package{}
		for _, path := range paths {
			pkg, err := check(path)
			if err != nil {
				return nil, err
			}
			loaded[path] = pkg
		}
		return loaded, nil
	}
}

type importerFunc func(path string) (*types.Package, error)

func (f importerFunc) Import(path string) (*types.Package, error) { return f(path) }

func (f importerFunc) ImportFrom(path, _ string, _ types.ImportMode) (*types.Package, error) {
	return f(path)
}

const (
	basePath  = testModule + "/base"
	quotePath = testModule + "/quote"
	restyPath = testModule + "/resty"
	localPath = testModule + "/local"
)

// middlewareSources declare quote.Service, which embeds base.Pinger, and the
// resty client that implements it.
var middlewareSources = map[string]string{
	basePath: `package base

import "context"

type Pinger interface {
	Ping(ctx context.Context) error
}
`,
	quotePath: `package quote

import (
	"context"

	"example.com/app/base"
)

type Price struct{ Value float64 }

type Service interface {
	base.Pinger
	Quote(ctx context.Context, symbol string) (Price, error)
	Name(ctx context.Context) string
	Tags(ctx context.Context, prefix string, extra ...string) ([]string, error)
	Next(ctx context.Context) error
}
`,
	restyPath: `package resty

import "example.com/app/quote"

func New() quote.Service { return nil }

func NewNamed() quote.Service { return nil }
`,
}

var quoteImports = map[string]string{"quote": quotePath, "echo": echoPath}

// middlewareEntries add to validEntries the decorated resty.New and an
// order handler of app api that consumes quote.Service.
func middlewareEntries(anns ...annotation.Annotation) []annotation.Entry {
	if len(anns) == 0 {
		anns = []annotation.Annotation{ann("Provide", "index=0"), ann("Middleware", "log, hystrix")}
	}
	entries := []annotation.Entry{
		decl("resty", "", "New", nil, []string{"quote.Service"}, anns...),
		decl("quoting", "", "NewHandler", []string{"quote.Service"}, []string{"*Handler"}, ann("Provide", "index=0")),
		decl("quoting", "Handler", "Get", echoHandler, errResult, ann("REST", "method=GET, path=/quotes/:symbol")),
	}
	for i := range entries {
		entries[i].Imports = quoteImports
	}
	return append(validEntries(), entries...)
}

func newMiddlewareSpec(t *testing.T, entries []annotation.Entry, sources map[string]string) (*Spec, error) {
	t.Helper()
	g, err := inject.NewGraphFromEntries(context.Background(), entries)
	require.NoError(t, err)
	return NewSpec(entries, g, testLoader(t, sources))
}

func TestNewSpec_Middleware(t *testing.T) {
	spec, err := newMiddlewareSpec(t, middlewareEntries(), middlewareSources)
	require.NoError(t, err)

	api, worker := spec.Apps[0], spec.Apps[1]
	require.Len(t, api.Middlewares, 1, "api serves the handler that consumes quote.Service")
	m := api.Middlewares[0]
	assert.Equal(t, restyPath+".New", m.Provider.Key())
	assert.Equal(t, 0, m.Index)
	assert.Equal(t, []string{"log", "hystrix"}, m.Plugins, "in the order written: the first is the outermost")
	require.NotNil(t, m.Interface)
	assert.Equal(t, quotePath+".Service", m.Interface.Obj().Pkg().Path()+"."+m.Interface.Obj().Name())

	assert.Empty(t, worker.Middlewares, "worker never builds resty.New")
}

func TestNewSpec_MiddlewareWithoutLoader(t *testing.T) {
	spec, err := newTestSpec(t, validEntries())
	require.NoError(t, err, "without @Middleware no loader is needed")
	assert.Empty(t, spec.Apps[0].Middlewares)

	g, err := inject.NewGraphFromEntries(context.Background(), middlewareEntries())
	require.NoError(t, err)
	_, err = NewSpec(middlewareEntries(), g, nil)
	require.Error(t, err)
	assert.True(t, errors.IsNotValid(err), "unexpected error kind: %v", err)
}

func TestNewSpec_MiddlewareErrors(t *testing.T) {
	provide := ann("Provide", "index=0")
	middleware := func(v string) annotation.Annotation { return ann("Middleware", v) }

	// quote declares Service as body and New returning it
	quote := func(body string) map[string]string {
		return map[string]string{quotePath: "package quote\n\nimport \"context\"\n\nvar _ context.Context\n\n" + body + "\n\nfunc New() Service { return nil }\n"}
	}
	newQuote := func(anns ...annotation.Annotation) []annotation.Entry {
		return []annotation.Entry{decl("quote", "", "New", nil, []string{"Service"}, anns...)}
	}
	method := func(m string) map[string]string { return quote("type Service interface {\n" + m + "\n}") }
	get := method("Get(ctx context.Context) error")

	testCases := []struct {
		name    string
		entries []annotation.Entry
		sources map[string]string
		kind    func(error) bool
		msg     string
	}{
		{"empty list", newQuote(provide, middleware("")), get, errors.IsNotValid, "at least one plugin"},
		{"duplicate plugin", newQuote(provide, middleware("log, log")), get, errors.IsNotValid, "log twice"},
		{"unknown plugin", newQuote(provide, middleware("log, cache")), get, errors.IsNotValid, "cache"},
		{"key instead of names", newQuote(provide, middleware("plugins=log")), get, errors.IsNotValid, "plugin names"},
		{"repeated annotation", newQuote(provide, middleware("log"), middleware("hystrix")), get, errors.IsNotValid, "one @Middleware"},
		{"without provide", newQuote(middleware("log")), get, errors.IsNotValid, "@Provide"},
		{"config struct", append(newQuote(provide), decl("quote", "Options", "", nil, nil, ann("Config", ""), middleware("log"))), get,
			errors.IsNotValid, "only a @Provide function"},
		{"group provide", newQuote(ann("Provide", "group=quotes, index=0"), middleware("log")), get, errors.IsNotValid, "group"},
		{"concrete result", []annotation.Entry{decl("quote", "", "New", nil, []string{"*Service"}, provide, middleware("log"))},
			map[string]string{quotePath: "package quote\n\nimport \"context\"\n\ntype Service struct{}\n\nfunc (*Service) Get(ctx context.Context) error { return nil }\n\nfunc New() *Service { return nil }\n"},
			errors.IsNotValid, "named interface"},
		{"generic interface", []annotation.Entry{decl("quote", "", "New", nil, []string{"Service[int]"}, provide, middleware("log"))},
			map[string]string{quotePath: "package quote\n\nimport \"context\"\n\ntype Service[T any] interface {\nGet(ctx context.Context) (T, error)\n}\n\nfunc New() Service[int] { return nil }\n"},
			errors.IsNotValid, "generic"},
		{"interface without methods", newQuote(provide, middleware("log")), quote("type Service interface{}"), errors.IsNotValid, "no methods"},
		{"unexported method", newQuote(provide, middleware("log")), method("get(ctx context.Context) error"), errors.IsNotValid, "unexported"},
		{"method without context", newQuote(provide, middleware("log")), method("Get(id string) error"), errors.IsNotValid, "context.Context"},
		{"method without results", newQuote(provide, middleware("log")), method("Do(ctx context.Context)"), errors.IsNotValid, "error, (R, error) or R"},
		{"two results without error", newQuote(provide, middleware("log")), method("Get(ctx context.Context) (int, string)"), errors.IsNotValid, "error, (R, error) or R"},
		{"three results", newQuote(provide, middleware("log")), method("Get(ctx context.Context) (int, string, error)"), errors.IsNotValid, "error, (R, error) or R"},
		{"unexported parameter type", newQuote(provide, middleware("log")),
			quote("type key string\n\ntype Service interface {\nGet(ctx context.Context, k key) error\n}"), errors.IsNotValid, "unexported type"},
		{"unexported type in a result", newQuote(provide, middleware("log")),
			quote("type item struct{}\n\ntype Service interface {\nList(ctx context.Context) ([]*item, error)\n}"), errors.IsNotValid, "unexported type"},
		{"same interface decorated twice", append(newQuote(provide, middleware("log")),
			decl("quote", "", "NewOther", nil, []string{"Service"}, ann("Provide", "name=other, index=0"), middleware("hystrix"))),
			map[string]string{quotePath: "package quote\n\nimport \"context\"\n\ntype Service interface {\nGet(ctx context.Context) error\n}\n\nfunc New() Service { return nil }\n\nfunc NewOther() Service { return nil }\n"},
			errors.IsAlreadyExists, "quote.service.get"},
		{"app named middleware", []annotation.Entry{decl("probe", "", "Live", echoHandler, errResult, ann("REST", "method=GET, path=/live, app=middleware"))},
			nil, errors.IsNotValid, "app name"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := newMiddlewareSpec(t, tc.entries, tc.sources)
			require.Error(t, err)
			assert.True(t, tc.kind(err), "unexpected error kind: %v", err)
			assert.Contains(t, err.Error(), tc.msg)
		})
	}
}

func TestLoadPackages(t *testing.T) {
	const orderPath = "github.com/xgodev/boost/builder/examples/simple/order"
	pkgs, err := LoadPackages(".")(orderPath)
	require.NoError(t, err)
	require.Contains(t, pkgs, orderPath)
	assert.NotNil(t, pkgs[orderPath].Scope().Lookup("NewService"))

	_, err = LoadPackages(".")("github.com/xgodev/boost/builder/examples/simple/missing")
	require.Error(t, err)
}

func renderMiddleware(t *testing.T, entries []annotation.Entry, sources map[string]string) map[string][]byte {
	t.Helper()
	spec, err := newMiddlewareSpec(t, entries, sources)
	require.NoError(t, err)
	files, err := NewGenerator(spec, testOptions("/src/app")).Render(context.Background())
	require.NoError(t, err)
	for name, src := range files {
		assertImportsUsed(t, name, src)
	}
	return files
}

func TestGenerator_RenderMiddleware(t *testing.T) {
	files := renderMiddleware(t, middlewareEntries(), middlewareSources)

	const decorator = "builder/middleware/resty/new_middleware.go"
	require.Contains(t, fileNames(files), decorator)

	golden := filepath.Join("testdata", "golden", "middleware", "new_middleware.go.golden")
	if *update {
		require.NoError(t, os.MkdirAll(filepath.Dir(golden), 0o755))
		require.NoError(t, os.WriteFile(golden, files[decorator], 0o644))
	}
	want, err := os.ReadFile(golden)
	require.NoError(t, err, "run go test ./builder/ -run TestGenerator_RenderMiddleware -update")
	assert.Equal(t, string(want), string(files[decorator]))

	base, err := testOptions("/src/app").BaseImport()
	require.NoError(t, err)
	decoratorPath := base + "/" + inject.GenPathUnder(base, "builder/middleware", restyPath)
	call := inject.Alias("mw", "resty", decoratorPath) + ".DecorateNew()"

	api := string(files["builder/api/modules.go"])
	assert.Contains(t, api, `"`+decoratorPath+`"`)
	assert.Contains(t, api, call+",\n\t\tmultiserverfx.Module()", "the decorators come after the doors, before the multiserver")
	assert.NotContains(t, string(files["builder/worker/modules.go"]), call, "worker never builds resty.New")
}

func TestGenerator_RenderNamedMiddleware(t *testing.T) {
	entries := middlewareEntries(ann("Provide", "name=x, index=0"), ann("Middleware", "log"))
	for i := range entries {
		if entries[i].Func.Name == "NewHandler" {
			entries[i].Annotations = append(entries[i].Annotations, ann("Inject", "name=x, index=0"))
		}
	}
	files := renderMiddleware(t, entries, middlewareSources)

	src := string(files["builder/middleware/resty/new_middleware.go"])
	assert.Contains(t, src, "fx.Decorate(fx.Annotate(decorateNew, fx.ParamTags(``, `name:\"x\"`), fx.ResultTags(`name:\"x\"`)))",
		"the decorated value keeps its name")
}

// implementationEntries make resty.New (@Default) and local.New the
// implementations of quote.Service, each with its own @Middleware.
func implementationEntries() ([]annotation.Entry, map[string]string) {
	entries := middlewareEntries(ann("Provide", "index=0"), ann("Middleware", "log, hystrix"), ann("Default", ""))
	local := decl("local", "", "New", nil, []string{"quote.Service"}, ann("Provide", "index=0"), ann("Middleware", "log"))
	local.Imports = quoteImports

	sources := map[string]string{localPath: "package local\n\nimport \"example.com/app/quote\"\n\nfunc New() quote.Service { return nil }\n"}
	for path, src := range middlewareSources {
		sources[path] = src
	}
	return append(entries, local), sources
}

func TestNewSpec_MiddlewareOnImplementation(t *testing.T) {
	entries, sources := implementationEntries()
	spec, err := newMiddlewareSpec(t, entries, sources)
	require.NoError(t, err, "two implementations may decorate the same interface: only one is wired")

	api := spec.Apps[0]
	require.Len(t, api.Middlewares, 2, "the consumer reaches every implementation")
	for _, m := range api.Middlewares {
		require.NotNil(t, m.Implementation, m.Provider.Key())
		assert.Equal(t, "resty", m.Implementation.Default)
		assert.Equal(t, []string{"local", "resty"}, m.Implementation.Names())
	}
}

func TestGenerator_RenderMiddlewareOnImplementation(t *testing.T) {
	entries, sources := implementationEntries()
	files := renderMiddleware(t, entries, sources)

	base, err := testOptions("/src/app").BaseImport()
	require.NoError(t, err)
	switchPath := base + "/" + inject.GenPathUnder(base, "implementation", quotePath)
	chosen := inject.Alias("impl", "quote", switchPath) + ".ServiceImplementation()"

	for member, path := range map[string]string{"resty": restyPath, "local": localPath} {
		src := string(files[inject.GenPathUnder(base, "builder/middleware", path)+"/new_middleware.go"])
		assert.Contains(t, src, `"`+switchPath+`"`)
		assert.Contains(t, src, "if "+chosen+" != \""+member+"\" {\n\t\treturn fx.Options()\n\t}",
			"the decorator of %s runs only when %s is the chosen implementation", member, member)
	}

	modules := string(files["builder/api/modules.go"])
	for _, path := range []string{restyPath, localPath} {
		assert.Contains(t, modules, inject.Alias("mw", filepath.Base(path), base+"/"+inject.GenPathUnder(base, "builder/middleware", path))+".DecorateNew()")
	}
}
