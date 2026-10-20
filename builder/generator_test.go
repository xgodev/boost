package builder

import (
	"context"
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost/annotation"
	"github.com/xgodev/boost/inject"
	"github.com/xgodev/boost/model/errors"
)

var update = flag.Bool("update", false, "rewrite the golden files of the generator tests")

func testOptions(root string) Options {
	return Options{
		GeneratorOptions: inject.GeneratorOptions{ModulePath: testModule, ModuleDir: root, OutputDir: filepath.Join(root, "gen")},
		CmdDir:           filepath.Join(root, "cmd"),
	}
}

func render(t *testing.T, entries []annotation.Entry) (map[string][]byte, error) {
	t.Helper()
	spec, err := newTestSpec(t, entries)
	require.NoError(t, err)
	return NewGenerator(spec, testOptions("/src/app")).Render(context.Background())
}

func fileNames(files map[string][]byte) []string {
	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func TestGenerator_Render(t *testing.T) {
	files, err := render(t, validEntries())
	require.NoError(t, err)

	assert.Equal(t, []string{
		"builder/api/admin_module.go",
		"builder/api/grpc_module.go",
		"builder/api/mcp_module.go",
		"builder/api/modules.go",
		"builder/api/ops_module.go",
		"builder/api/rest_module.go",
		"builder/worker/function_module.go",
		"builder/worker/modules.go",
	}, fileNames(files))

	for name, src := range files {
		assertImportsUsed(t, name, src)

		golden := filepath.Join("testdata", "golden", filepath.FromSlash(strings.TrimPrefix(name, "builder/"))+".golden")
		if *update {
			require.NoError(t, os.MkdirAll(filepath.Dir(golden), 0o755))
			require.NoError(t, os.WriteFile(golden, src, 0o644))
			continue
		}
		want, err := os.ReadFile(golden)
		require.NoError(t, err, "run go test ./builder/ -run TestGenerator_Render -update")
		assert.Equal(t, string(want), string(src), name)
	}
}

// assertImportsUsed fails when the file imports a package it never refers
// to, which would not compile.
func assertImportsUsed(t *testing.T, name string, src []byte) {
	t.Helper()
	file, err := parser.ParseFile(token.NewFileSet(), name, src, 0)
	require.NoError(t, err, name)

	used := map[string]bool{}
	ast.Inspect(file, func(n ast.Node) bool {
		if sel, ok := n.(*ast.SelectorExpr); ok {
			if x, ok := sel.X.(*ast.Ident); ok {
				used[x.Name] = true
			}
		}
		return true
	})
	for _, imp := range file.Imports {
		alias := path.Base(strings.Trim(imp.Path.Value, `"`))
		if imp.Name != nil {
			alias = imp.Name.Name
		}
		assert.True(t, used[alias], "%s imports %s without using it", name, imp.Path.Value)
	}
}

func TestGenerator_RenderIsStable(t *testing.T) {
	first, err := render(t, validEntries())
	require.NoError(t, err)

	reversed := validEntries()
	for i, j := 0, len(reversed)-1; i < j; i, j = i+1, j-1 {
		reversed[i], reversed[j] = reversed[j], reversed[i]
	}
	second, err := render(t, reversed)
	require.NoError(t, err)

	assert.Equal(t, first, second, "the order of the annotations does not change the output")
}

// TestGenerator_RenderKebabApp checks an app named in kebab case, as deploys
// are (#242): cmd/<app> and gen/builder/<app> keep the hyphens, the Go
// package drops them and the config path element is camelCase, the key an
// environment variable reaches (ORDERS__UPDATER_FUNCTION_ADAPTER is
// ordersUpdater.function.adapter).
func TestGenerator_RenderKebabApp(t *testing.T) {
	entries := []annotation.Entry{
		decl("probe", "", "Live", echoHandler, errResult, ann("REST", "method=GET, path=/live, server=admin, app=orders-updater")),
		decl("events", "", "New", nil, []string{"Handler"}, ann("Provide", "index=0")),
		decl("events", "Handler", "Handle", eventParams, eventResults, ann("Function", "adapter=[pubsub, nats], app=orders-updater")),
	}
	spec, err := newTestSpec(t, entries)
	require.NoError(t, err)
	require.Len(t, spec.Apps, 1)
	assert.Equal(t, "orders-updater", spec.Apps[0].Name)

	gen := NewGenerator(spec, testOptions("/src/app"))
	files, err := gen.Render(context.Background())
	require.NoError(t, err)
	mains, err := gen.RenderCmd(context.Background())
	require.NoError(t, err)
	for name, src := range mains {
		files["cmd/"+name] = src
	}

	assert.Equal(t, []string{
		"builder/orders-updater/admin_module.go",
		"builder/orders-updater/function_module.go",
		"builder/orders-updater/modules.go",
		"cmd/orders-updater/main.go",
	}, fileNames(files))

	for name, src := range files {
		assertImportsUsed(t, name, src)
		if strings.HasPrefix(name, "builder/") {
			assert.Contains(t, string(src), "\npackage ordersupdater\n", name)
		}

		golden := filepath.Join("testdata", "golden", "kebab", path.Base(name)+".golden")
		if *update {
			require.NoError(t, os.MkdirAll(filepath.Dir(golden), 0o755))
			require.NoError(t, os.WriteFile(golden, src, 0o644))
			continue
		}
		want, err := os.ReadFile(golden)
		require.NoError(t, err, "run go test ./builder/ -run TestGenerator_RenderKebabApp -update")
		assert.Equal(t, string(want), string(src), name)
	}
	assert.Contains(t, string(files["builder/orders-updater/admin_module.go"]), `"ordersUpdater.server.admin"`)
	assert.Contains(t, string(files["builder/orders-updater/function_module.go"]), `"ordersUpdater.function.adapter"`)
	assert.Contains(t, string(files["cmd/orders-updater/main.go"]), `app "example.com/app/gen/builder/orders-updater"`)
}

// TestConfigName checks that the config path element of an app is the key
// boost's env loader reads from its environment variable: hyphens become
// camelCase, as ORDERS__UPDATER_X reaches ordersUpdater.x.
// TestGenerator_RenderJob renders app order-sweeper, whose @Job is a method,
// and app cleanup, whose @Job is a package function: each has a job door
// and the main every app has.
func TestGenerator_RenderJob(t *testing.T) {
	entries := append(jobEntries("app=order-sweeper"),
		decl("cleanup", "", "Run", []string{"context.Context"}, errResult, ann("Job", "app=cleanup")))
	spec, err := newTestSpec(t, entries)
	require.NoError(t, err)

	gen := NewGenerator(spec, testOptions("/src/app"))
	files, err := gen.Render(context.Background())
	require.NoError(t, err)
	mains, err := gen.RenderCmd(context.Background())
	require.NoError(t, err)
	for name, src := range mains {
		files["cmd/"+name] = src
	}

	assert.Equal(t, []string{
		"builder/cleanup/job_module.go",
		"builder/cleanup/modules.go",
		"builder/order-sweeper/job_module.go",
		"builder/order-sweeper/modules.go",
		"cmd/cleanup/main.go",
		"cmd/order-sweeper/main.go",
	}, fileNames(files))

	for name, src := range files {
		assertImportsUsed(t, name, src)

		golden := filepath.Join("testdata", "golden", "job", filepath.FromSlash(name)+".golden")
		if *update {
			require.NoError(t, os.MkdirAll(filepath.Dir(golden), 0o755))
			require.NoError(t, os.WriteFile(golden, src, 0o644))
			continue
		}
		want, err := os.ReadFile(golden)
		require.NoError(t, err, "run go test ./builder/ -run TestGenerator_RenderJob -update")
		assert.Equal(t, string(want), string(src), name)
	}

	door := string(files["builder/order-sweeper/job_module.go"])
	assert.Contains(t, door, "jobfx.Module()")
	assert.Contains(t, door, "return h0.Run")
	assert.NotContains(t, string(files["builder/cleanup/job_module.go"]), "func(h0", "a package function needs no receiver")
	assert.NotContains(t, string(files["builder/order-sweeper/modules.go"]), "multiserver", "a job app has no server")
}

func TestConfigName(t *testing.T) {
	for app, want := range map[string]string{
		"api":            "api",
		"orders-updater": "ordersUpdater",
		"orders-v2-api":  "ordersV2Api",
		"orders-2":       "orders2",
	} {
		assert.Equal(t, want, configName(app), app)
	}
}

// TestGenerator_RenderFunctionAdapterList checks the function door of an
// adapter list: the key that chooses the adapter and the switch that wires
// only the chosen one. A single adapter renders as in TestGenerator_Render.
func TestGenerator_RenderFunctionAdapterList(t *testing.T) {
	files, err := render(t, []annotation.Entry{
		decl("events", "", "New", nil, []string{"Handler"}, ann("Provide", "index=0")),
		decl("events", "Handler", "Handle", eventParams, eventResults, ann("Function", "adapter=[pubsub, nats], app=worker")),
	})
	require.NoError(t, err)

	name := "builder/worker/function_module.go"
	src, ok := files[name]
	require.True(t, ok, "files: %v", fileNames(files))
	assertImportsUsed(t, name, src)

	golden := filepath.Join("testdata", "golden", "adapters", "function_module.go.golden")
	if *update {
		require.NoError(t, os.MkdirAll(filepath.Dir(golden), 0o755))
		require.NoError(t, os.WriteFile(golden, src, 0o644))
		return
	}
	want, err := os.ReadFile(golden)
	require.NoError(t, err, "run go test ./builder/ -run TestGenerator_RenderFunctionAdapterList -update")
	assert.Equal(t, string(want), string(src))
}

// TestGenerator_RenderFunctionAdapterClient checks that an adapter reading
// with a client inject builds takes the app's default client of that type,
// so the publisher and the @Inject consumers share it and fx gets it once:
// the module inject generates for it (Spec.DefaultClients), or the module of
// the inject constructor providing it.
func TestGenerator_RenderFunctionAdapterClient(t *testing.T) {
	newHandler := decl("events", "", "New", nil, []string{"Handler"}, ann("Provide", "index=0"))
	handle := decl("events", "Handler", "Handle", eventParams, eventResults, ann("Function", "adapter=[nats, cenats, pubsub], app=worker"))
	const name = "builder/worker/function_module.go"

	t.Run("inject builds the client", func(t *testing.T) {
		spec, err := newTestSpec(t, []annotation.Entry{newHandler, handle})
		require.NoError(t, err)
		assert.Equal(t, []string{"cloud.google.com/go/pubsub/v2.Client", "github.com/nats-io/nats.go.Conn"}, spec.DefaultClients())

		files, err := NewGenerator(spec, testOptions("/src/app")).Render(context.Background())
		require.NoError(t, err)
		src := files[name]
		assertImportsUsed(t, name, src)

		golden := filepath.Join("testdata", "golden", "adapters", "client_function_module.go.golden")
		if *update {
			require.NoError(t, os.WriteFile(golden, src, 0o644))
			return
		}
		want, err := os.ReadFile(golden)
		require.NoError(t, err, "run go test ./builder/ -run TestGenerator_RenderFunctionAdapterClient -update")
		assert.Equal(t, string(want), string(src))
	})

	t.Run("the app provides the client", func(t *testing.T) {
		conn := decl("broker", "", "NewConn", nil, []string{"*nats.Conn"}, ann("Provide", "index=0"))
		conn.Imports = map[string]string{"nats": "github.com/nats-io/nats.go"}
		spec, err := newTestSpec(t, []annotation.Entry{newHandler, handle, conn})
		require.NoError(t, err)
		assert.Equal(t, []string{"cloud.google.com/go/pubsub/v2.Client"}, spec.DefaultClients(),
			"inject builds no *nats.Conn the app provides")

		files, err := NewGenerator(spec, testOptions("/src/app")).Render(context.Background())
		require.NoError(t, err)
		src := string(files[name])
		assertImportsUsed(t, name, files[name])

		_, call := inject.ModuleOf(testModule+"/gen", conn)
		assert.Equal(t, 2, strings.Count(src, call+","), "nats and cenats take the app's *nats.Conn:\n%s", src)
		assert.NotContains(t, src, "DefaultConnModule()")
		assert.Contains(t, src, "DefaultClientModule()")
	})
}

// TestGenerator_RenderFunctionNamedClient checks @Function(client=prices):
// the adapters reading with a client inject builds take the client named
// prices, the module inject generates for it (Spec.NamedClients) or the one
// of the inject constructor providing it, and no default client (#208).
func TestGenerator_RenderFunctionNamedClient(t *testing.T) {
	newHandler := decl("events", "", "New", nil, []string{"Handler"}, ann("Provide", "index=0"))
	handle := decl("events", "Handler", "Handle", eventParams, eventResults, ann("Function", "adapter=[nats, cenats, pubsub], app=worker, client=prices"))
	const name = "builder/worker/function_module.go"

	t.Run("inject builds the client", func(t *testing.T) {
		spec, err := newTestSpec(t, []annotation.Entry{newHandler, handle})
		require.NoError(t, err)
		assert.Empty(t, spec.DefaultClients(), "the adapters read with the prices clients")
		assert.Equal(t, []inject.ClientRef{
			{ID: "cloud.google.com/go/pubsub/v2.Client", Name: "prices"},
			{ID: "github.com/nats-io/nats.go.Conn", Name: "prices"},
		}, spec.NamedClients())

		files, err := NewGenerator(spec, testOptions("/src/app")).Render(context.Background())
		require.NoError(t, err)
		src := files[name]
		assertImportsUsed(t, name, src)

		golden := filepath.Join("testdata", "golden", "adapters", "named_client_function_module.go.golden")
		if *update {
			require.NoError(t, os.WriteFile(golden, src, 0o644))
			return
		}
		want, err := os.ReadFile(golden)
		require.NoError(t, err, "run go test ./builder/ -run TestGenerator_RenderFunctionNamedClient -update")
		assert.Equal(t, string(want), string(src))
	})

	t.Run("the app provides the client", func(t *testing.T) {
		conn := decl("broker", "", "NewConn", nil, []string{"*nats.Conn"}, ann("Provide", "index=0, name=prices"))
		conn.Imports = map[string]string{"nats": "github.com/nats-io/nats.go"}
		spec, err := newTestSpec(t, []annotation.Entry{newHandler, handle, conn})
		require.NoError(t, err)
		assert.Equal(t, []inject.ClientRef{{ID: "cloud.google.com/go/pubsub/v2.Client", Name: "prices"}}, spec.NamedClients(),
			"inject builds no *nats.Conn the app provides")

		files, err := NewGenerator(spec, testOptions("/src/app")).Render(context.Background())
		require.NoError(t, err)
		src := string(files[name])
		assertImportsUsed(t, name, files[name])

		_, call := inject.ModuleOf(testModule+"/gen", conn)
		assert.Equal(t, 2, strings.Count(src, call+","), "nats and cenats take the app's prices *nats.Conn:\n%s", src)
		assert.Equal(t, 3, strings.Count(src, `ModuleWithNamedClient[*src_cloudevents_c136a3fd.Event]("prices")`), "nats, cenats and pubsub:\n%s", src)
		assert.NotContains(t, src, "PricesConnModule()")
		assert.Contains(t, src, "PricesClientModule()")
	})
}

func TestGenerator_RenderErrors(t *testing.T) {
	handle := func(results ...string) annotation.Entry {
		return decl("events", "", "Handle", eventParams, results, ann("Function", "adapter=nats"))
	}
	grpcServer := func(result string) annotation.Entry {
		return decl("health", "", "NewServer", nil, []string{result},
			ann("Provide", "index=0"), ann("GRPC", "register=healthpb.RegisterHealthServer"))
	}

	tests := []struct {
		name    string
		entries []annotation.Entry
		want    string
	}{
		{"function result from a package the file does not import", []annotation.Entry{handle("*missing.Event", "error")}, "missing.Event"},
		{"function result of an unexported type", []annotation.Entry{handle("*result", "error")}, "result"},
		{"gRPC value of an unexported type", []annotation.Entry{grpcServer("*server")}, "server"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := render(t, tt.entries)
			require.Error(t, err)
			assert.True(t, errors.IsNotValid(err), "got %v", err)
			assert.Contains(t, err.Error(), tt.want)
		})
	}
}

func TestGenerator_RenderNeedsOutputInsideTheModule(t *testing.T) {
	spec, err := newTestSpec(t, validEntries())
	require.NoError(t, err)

	opts := Options{GeneratorOptions: inject.GeneratorOptions{ModulePath: testModule, ModuleDir: "/src/app", OutputDir: "/elsewhere/gen"}}
	_, err = NewGenerator(spec, opts).Render(context.Background())
	assert.True(t, errors.IsNotValid(err), "got %v", err)
}

func TestGenerator_RenderCmd(t *testing.T) {
	spec, err := newTestSpec(t, validEntries())
	require.NoError(t, err)
	files, err := NewGenerator(spec, testOptions("/src/app")).RenderCmd(context.Background())
	require.NoError(t, err)

	assert.Equal(t, []string{"api/main.go", "worker/main.go"}, fileNames(files))

	for name, src := range files {
		assertImportsUsed(t, name, src)

		golden := filepath.Join("testdata", "golden", "cmd", filepath.FromSlash(name)+".golden")
		if *update {
			require.NoError(t, os.MkdirAll(filepath.Dir(golden), 0o755))
			require.NoError(t, os.WriteFile(golden, src, 0o644))
			continue
		}
		want, err := os.ReadFile(golden)
		require.NoError(t, err, "run go test ./builder/ -run TestGenerator_RenderCmd -update")
		assert.Equal(t, string(want), string(src), name)
	}
}

func TestGenerator_RenderCmdWithoutCmdDir(t *testing.T) {
	spec, err := newTestSpec(t, validEntries())
	require.NoError(t, err)
	opts := testOptions("/src/app")
	opts.CmdDir = ""

	files, err := NewGenerator(spec, opts).RenderCmd(context.Background())
	require.NoError(t, err)
	assert.Empty(t, files, "no cmd dir, no main")
}

func TestGenerator_RenderCmdNeedsCmdInsideTheModule(t *testing.T) {
	spec, err := newTestSpec(t, validEntries())
	require.NoError(t, err)
	opts := testOptions("/src/app")
	opts.CmdDir = "/elsewhere/cmd"

	_, err = NewGenerator(spec, opts).RenderCmd(context.Background())
	assert.True(t, errors.IsNotValid(err), "got %v", err)
}

func TestGenerator_Generate(t *testing.T) {
	root := t.TempDir()
	opts := testOptions(root)
	write := func(dir, rel, body string) {
		p := filepath.Join(dir, filepath.FromSlash(rel))
		require.NoError(t, os.MkdirAll(filepath.Dir(p), 0o755))
		require.NoError(t, os.WriteFile(p, []byte(body), 0o644))
	}
	write(opts.OutputDir, "builder/old/rest_module.go", generatedHeader+"\n\npackage old\n")
	write(opts.OutputDir, "builder/api/handwritten.go", "package api\n")
	write(opts.OutputDir, "module.go", "// Code generated by inject; DO NOT EDIT.\n\npackage gen\n")
	write(opts.CmdDir, "old/main.go", generatedHeader+"\n\npackage main\n")
	write(opts.CmdDir, "tool/main.go", "package main\n")

	spec, err := newTestSpec(t, validEntries())
	require.NoError(t, err)
	gen := NewGenerator(spec, opts)
	require.NoError(t, gen.Generate(context.Background()))

	files, err := gen.Render(context.Background())
	require.NoError(t, err)
	for name, want := range files {
		got, err := os.ReadFile(filepath.Join(opts.OutputDir, filepath.FromSlash(name)))
		require.NoError(t, err)
		assert.Equal(t, string(want), string(got), name)
	}

	assert.NoFileExists(t, filepath.Join(opts.OutputDir, "builder", "old", "rest_module.go"), "stale builder files are removed")
	assert.NoDirExists(t, filepath.Join(opts.OutputDir, "builder", "old"))
	assert.FileExists(t, filepath.Join(opts.OutputDir, "builder", "api", "handwritten.go"), "handwritten files are kept")
	assert.FileExists(t, filepath.Join(opts.OutputDir, "module.go"), "inject files are not builder's")

	mains, err := gen.RenderCmd(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, mains)
	for name, want := range mains {
		got, err := os.ReadFile(filepath.Join(opts.CmdDir, filepath.FromSlash(name)))
		require.NoError(t, err)
		assert.Equal(t, string(want), string(got), name)
	}
	assert.NoDirExists(t, filepath.Join(opts.CmdDir, "old"), "the main of a removed app is removed")
	assert.FileExists(t, filepath.Join(opts.CmdDir, "tool", "main.go"), "handwritten commands are kept")
}

func TestGenerator_GenerateRefusesHandwrittenMain(t *testing.T) {
	root := t.TempDir()
	opts := testOptions(root)
	mine := filepath.Join(opts.CmdDir, "api", "main.go")
	require.NoError(t, os.MkdirAll(filepath.Dir(mine), 0o755))
	require.NoError(t, os.WriteFile(mine, []byte("package main\n\n// mine\n"), 0o644))

	spec, err := newTestSpec(t, validEntries())
	require.NoError(t, err)
	err = NewGenerator(spec, opts).Generate(context.Background())
	require.Error(t, err)
	assert.True(t, errors.IsAlreadyExists(err), "got %v", err)
	assert.Contains(t, err.Error(), "main.go")

	got, err := os.ReadFile(mine)
	require.NoError(t, err)
	assert.Equal(t, "package main\n\n// mine\n", string(got), "the handwritten main is not overwritten")
	assert.NoDirExists(t, opts.OutputDir, "nothing is written")
	assert.NoFileExists(t, filepath.Join(opts.CmdDir, "worker", "main.go"), "nothing is written")
}

func TestImports_Qualify(t *testing.T) {
	entry := decl("orders", "", "New", nil, nil)
	src := inject.Alias("src", "orders", testModule+"/orders")
	ev := inject.Alias("src", "cloudevents", "github.com/cloudevents/sdk-go/v2")

	tests := []struct{ expr, want string }{
		{"error", "error"},
		{"*Handler", "*" + src + ".Handler"},
		{"[]map[string]*cloudevents.Event", "[]map[string]*" + ev + ".Event"},
		{"func(ctx context.Context, h Handler) error", "func(ctx " + inject.Alias("src", "context", "context") + ".Context, h " + src + ".Handler) error"},
		{"struct{ Name string }", "struct{Name string}"},
		{"interface{ Do() error }", "interface{Do() error}"},
		{"Box[int]", src + ".Box[int]"},
		{"chan<- any", "chan<- any"},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			got, err := newImports().qualify(entry, tt.expr)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}

	imports := newImports()
	_, err := imports.qualify(entry, "*Handler")
	require.NoError(t, err)
	_, err = imports.qualify(entry, "orders.Item") // same path through the file import
	require.NoError(t, err)
	assert.Equal(t, []inject.ImportData{{Alias: src, Path: path.Join(testModule, "orders")}}, imports.sorted(),
		"a package is imported once whatever name the file gives it")
}

func TestGenerator_RenderIncludesPluginModules(t *testing.T) {
	plugins := pluginEntries()
	files, err := render(t, append(validEntries(), plugins...))
	require.NoError(t, err)

	base, err := testOptions("/src/app").BaseImport()
	require.NoError(t, err)
	imp, admin := inject.ModuleOf(base, plugins[0])
	_, rest := inject.ModuleOf(base, plugins[1])
	_, grpc := inject.ModuleOf(base, plugins[2])
	_, mcp := inject.ModuleOf(base, plugins[3])
	_, ops := inject.ModuleOf(base, plugins[4])

	for name, src := range files {
		assertImportsUsed(t, name, src)
	}
	assert.Contains(t, string(files["builder/api/admin_module.go"]), admin, "the admin door builds its plugins")
	assert.Contains(t, string(files["builder/api/admin_module.go"]), imp.Path)
	assert.Contains(t, string(files["builder/api/rest_module.go"]), rest)
	assert.NotContains(t, string(files["builder/api/rest_module.go"]), admin, "the admin plugins are admin's only")
	assert.Contains(t, string(files["builder/api/grpc_module.go"]), grpc)
	assert.Contains(t, string(files["builder/api/mcp_module.go"]), mcp)
	assert.Contains(t, string(files["builder/api/ops_module.go"]), ops)
	assert.NotContains(t, string(files["builder/api/mcp_module.go"]), ops, "the ops plugins are ops' only")
}

// TestGenerator_RenderDefaultServerNamedPlugins checks the rest, grpc and
// mcp doors of plugins provided under the name of their server: each door
// builds the module providing them and takes the named value.
func TestGenerator_RenderDefaultServerNamedPlugins(t *testing.T) {
	plugins := namedDefaultPluginEntries()
	files, err := render(t, append(validEntries(), plugins...))
	require.NoError(t, err)

	base, err := testOptions("/src/app").BaseImport()
	require.NoError(t, err)
	modules := map[string]string{}
	for i, door := range []string{"rest", "grpc", "mcp"} {
		_, modules[door] = inject.ModuleOf(base, plugins[i])
	}

	for _, door := range []string{"rest", "grpc", "mcp"} {
		name := "builder/api/" + door + "_module.go"
		src, ok := files[name]
		require.True(t, ok, "files: %v", fileNames(files))
		assertImportsUsed(t, name, src)
		for other, module := range modules {
			if other == door {
				assert.Contains(t, string(src), module, "the %s door builds its plugins", door)
			} else {
				assert.NotContains(t, string(src), module, "the %s plugins are %s's only", other, other)
			}
		}
		assert.NotContains(t, string(files["builder/api/admin_module.go"]), modules[door], "the %s plugins are %s's only", door, door)
		assert.NotContains(t, string(files["builder/api/ops_module.go"]), modules[door], "the %s plugins are %s's only", door, door)

		golden := filepath.Join("testdata", "golden", "plugins", door+"_module.go.golden")
		if *update {
			require.NoError(t, os.MkdirAll(filepath.Dir(golden), 0o755))
			require.NoError(t, os.WriteFile(golden, src, 0o644))
			continue
		}
		want, err := os.ReadFile(golden)
		require.NoError(t, err, "run go test ./builder/ -run TestGenerator_RenderDefaultServerNamedPlugins -update")
		assert.Equal(t, string(want), string(src), name)
	}
}

func TestGenerator_RenderIncludesMiddlewareProviders(t *testing.T) {
	drivers := driverEntries()
	files, err := render(t, append(drivers,
		decl("events", "", "Handle", eventParams, eventResults, ann("Function", "adapter=nats, app=worker, middlewares=[publisher]"))))
	require.NoError(t, err)

	base, err := testOptions("/src/app").BaseImport()
	require.NoError(t, err)
	imp, driver := inject.ModuleOf(base, drivers[0])
	_, named := inject.ModuleOf(base, drivers[1])

	for name, src := range files {
		assertImportsUsed(t, name, src)
	}
	src := string(files["builder/worker/function_module.go"])
	assert.Contains(t, src, driver, "the function door builds the driver the publisher middleware publishes with")
	assert.Contains(t, src, imp.Path)
	assert.NotContains(t, src, named, "a named driver is not the publisher's")
}
