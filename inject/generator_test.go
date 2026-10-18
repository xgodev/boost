package inject

import (
	"bytes"
	"context"
	"flag"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost/annotation"
	"github.com/xgodev/boost/extra/graph"
	"github.com/xgodev/boost/model/errors"
)

const testModule = "example.com/app"

// at builds the entry of a function declared in package pkg of testModule.
func at(pkg, name string, params, results []string, anns ...annotation.Annotation) annotation.Entry {
	entry := fn(name, params, results, anns...)
	entry.Path = testModule + "/" + pkg
	entry.Package = filepath.Base(pkg)
	return entry
}

func greetingEntries() []annotation.Entry {
	return []annotation.Entry{
		at("greet", "NewGreeter", nil, []string{"*Greeter"}, provide("index=0")),
		at("greet", "NewFormal", nil, []string{"*Greeter", "error"}, provide("name=formal,index=0")),
		at("names", "Alice", nil, []string{"string"}, provide("group=names,index=0")),
		at("app", "NewBoard", nil, []string{"*Board"}, provide("index=0")),
		at("app", "Run", []string{"*greet.Greeter", "*greet.Greeter", "[]string", "*Board"}, nil,
			inject("name=formal,index=0"), inject("group=names,index=2"), invoke()),
	}
}

func newTestGraph(t *testing.T, entries []annotation.Entry) *graph.Graph[Component] {
	t.Helper()
	g, err := NewGraphFromEntries(context.Background(), entries)
	require.NoError(t, err)
	return g
}

func testOptions(dir string) GeneratorOptions {
	return GeneratorOptions{ModulePath: testModule, ModuleDir: dir, OutputDir: filepath.Join(dir, "gen")}
}

func TestGenerator_Render(t *testing.T) {
	entries := greetingEntries()
	files, err := NewGenerator(newTestGraph(t, entries), testOptions("/mod")).Render(context.Background())
	require.NoError(t, err)

	assert.Equal(t, []string{
		"app/newboard_module.go",
		"app/run_module.go",
		"greet/newformal_module.go",
		"greet/newgreeter_module.go",
		"module.go",
		"names/alice_module.go",
	}, sortedFileNames(files))

	for name, src := range files {
		assert.True(t, bytes.HasPrefix(src, []byte(generatedHeader)), "%s has no generated header", name)
		_, err := parser.ParseFile(token.NewFileSet(), name, src, parser.AllErrors)
		assert.NoError(t, err, "%s is not valid Go", name)
	}

	run := string(files["app/run_module.go"])
	src, greet, names := srcAlias(entries[4]), genAlias(entries[0]), genAlias(entries[2])
	assert.Contains(t, run, "package app\n")
	assert.Contains(t, run, src+` "example.com/app/app"`)
	assert.Contains(t, run, greet+` "example.com/app/gen/greet"`)
	assert.Contains(t, run, names+` "example.com/app/gen/names"`)
	assert.NotContains(t, run, `"example.com/app/gen/app"`, "same package modules are not imported")
	assert.Contains(t, run, `fx.Module("example.com/app/app.Run",
			NewBoardModule(),
			`+greet+`.NewFormalModule(),
			`+greet+`.NewGreeterModule(),
			`+names+`.AliceModule(),
			fx.Invoke(fx.Annotate(`+src+".Run, fx.ParamTags(`name:\"formal\"`, \"\", `group:\"names\"`))),\n")

	assert.Contains(t, string(files["greet/newgreeter_module.go"]),
		"fx.Provide("+srcAlias(entries[0])+".NewGreeter)")
	assert.Contains(t, string(files["greet/newformal_module.go"]),
		"fx.Provide(fx.Annotate("+srcAlias(entries[1])+".NewFormal, fx.ResultTags(`name:\"formal\"`)))")
	assert.Contains(t, string(files["names/alice_module.go"]),
		"fx.Provide(fx.Annotate("+srcAlias(entries[2])+".Alice, fx.ResultTags(`group:\"names\"`)))")

	root := string(files["module.go"])
	assert.Contains(t, root, "package gen\n")
	for _, entry := range entries {
		assert.Contains(t, root, genAlias(entry)+"."+entry.Func.Name+"Module(),")
	}

	again, err := NewGenerator(newTestGraph(t, entries), testOptions("/mod")).Render(context.Background())
	require.NoError(t, err)
	assert.Equal(t, files, again, "render is deterministic")
}

func TestGenerator_RenderErrors(t *testing.T) {
	testCases := []struct {
		name    string
		entries []annotation.Entry
		opts    GeneratorOptions
	}{
		{
			name: "package main",
			entries: []annotation.Entry{func() annotation.Entry {
				e := at("cmd", "NewA", nil, []string{"*A"}, provide("index=0"))
				e.Package = "main"
				return e
			}()},
			opts: testOptions("/mod"),
		},
		{
			name:    "unexported function",
			entries: []annotation.Entry{at("p", "newA", nil, []string{"*A"}, provide("index=0"))},
			opts:    testOptions("/mod"),
		},
		{
			name:    "output outside the module",
			entries: greetingEntries(),
			opts:    GeneratorOptions{ModulePath: testModule, ModuleDir: "/mod", OutputDir: "/other/gen"},
		},
		{
			name:    "output is the module root",
			entries: greetingEntries(),
			opts:    GeneratorOptions{ModulePath: testModule, ModuleDir: "/mod", OutputDir: "/mod"},
		},
		{
			name:    "missing module path",
			entries: greetingEntries(),
			opts:    GeneratorOptions{ModuleDir: "/mod", OutputDir: "/mod/gen"},
		},
		{
			name: "two implementation sets with the same config key",
			entries: append(stockEntries()[:2], func() []annotation.Entry {
				var upper []annotation.Entry
				for i, pkg := range []string{"a", "b"} {
					e := at(pkg, "New", nil, []string{"Stock.Service"}, provide("index=0"))
					e.Imports = map[string]string{"Stock": testModule + "/Stock"}
					if i == 0 {
						e.Annotations = append(e.Annotations, deflt())
					}
					upper = append(upper, e)
				}
				return upper
			}()...),
			opts: testOptions("/mod"),
		},
		{
			name: "two functions generate the same file",
			entries: []annotation.Entry{
				at("p", "NewA", nil, []string{"*A"}, provide("index=0")),
				at("p", "Newa", nil, []string{"*B"}, provide("index=0")),
			},
			opts: testOptions("/mod"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewGenerator(newTestGraph(t, tc.entries), tc.opts).Render(context.Background())
			require.Error(t, err)
			assert.True(t, errors.IsNotValid(err) || errors.IsAlreadyExists(err), "unexpected error type: %v", err)
		})
	}
}

func TestGenerator_Generate(t *testing.T) {
	dir := t.TempDir()
	opts := testOptions(dir)

	write := func(rel, content string) {
		p := filepath.Join(opts.OutputDir, filepath.FromSlash(rel))
		require.NoError(t, os.MkdirAll(filepath.Dir(p), os.ModePerm))
		require.NoError(t, os.WriteFile(p, []byte(content), 0o644))
	}
	write("old/gone_module.go", generatedHeader+"\n\npackage old\n")
	write("handwritten.go", "package gen\n")
	write("notes.txt", generatedHeader)

	g := newTestGraph(t, greetingEntries())
	require.NoError(t, NewGenerator(g, opts).Generate(context.Background()))

	want, err := NewGenerator(g, opts).Render(context.Background())
	require.NoError(t, err)
	for name, src := range want {
		got, err := os.ReadFile(filepath.Join(opts.OutputDir, filepath.FromSlash(name)))
		require.NoError(t, err)
		assert.Equal(t, string(src), string(got), name)
	}

	assert.NoDirExists(t, filepath.Join(opts.OutputDir, "example.com/app/old"), "stale generated package is removed")
	assert.FileExists(t, filepath.Join(opts.OutputDir, "handwritten.go"), "files without the generated header are kept")
	assert.FileExists(t, filepath.Join(opts.OutputDir, "notes.txt"), "non-Go files are kept")

	// a component removed from the source disappears from the output
	require.NoError(t, NewGenerator(newTestGraph(t, greetingEntries()[:3]), opts).Generate(context.Background()))
	assert.NoDirExists(t, filepath.Join(opts.OutputDir, "example.com/app/app"))
	assert.FileExists(t, filepath.Join(opts.OutputDir, "names/alice_module.go"))
}

func TestWriteGraphviz(t *testing.T) {
	g := newTestGraph(t, []annotation.Entry{
		at("p", "NewA", nil, []string{"*A"}, provide("index=0")),
		at("p", "Run", []string{"*A"}, nil, invoke()),
	})

	var buf bytes.Buffer
	require.NoError(t, WriteGraphviz(&buf, g))
	assert.Equal(t, `digraph InjectGraph {
	graph [splines=true, overlap=false];
	node [fontname="Arial", style=filled];
	edge [color="#606060"];
	"example.com/app/p.NewA" [label="p.NewA", shape=ellipse, fillcolor="#90EE90"];
	"example.com/app/p.NewA" -> "example.com/app/p.Run";
	"example.com/app/p.Run" [label="p.Run", shape=octagon, fillcolor="#FFD700"];
}
`, buf.String())
}

func TestFindModule(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/app\n\ngo 1.22\n"), 0o644))
	sub := filepath.Join(dir, "internal", "pkg")
	require.NoError(t, os.MkdirAll(sub, os.ModePerm))

	modulePath, moduleDir, err := FindModule(sub)
	require.NoError(t, err)
	assert.Equal(t, "example.com/app", modulePath)
	assert.Equal(t, dir, moduleDir)

	broken := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(broken, "go.mod"), []byte("go 1.22\n"), 0o644))
	_, _, err = FindModule(broken)
	assert.True(t, errors.IsNotValid(err), "unexpected error: %v", err)
}

func TestTagList(t *testing.T) {
	two := 2
	assert.Equal(t, "", tagList(map[int]Annotation{0: {Index: &two}}))
	assert.Equal(t, `"", "", `+"`group:\"x\"`", tagList(map[int]Annotation{2: {Group: "x"}}))
	assert.Equal(t, "`name:\"a\"`", tagList(map[int]Annotation{0: {Name: "a"}}))
}

var update = flag.Bool("update", false, "rewrite the golden files of the generator tests")

// golden compares got with the file at name, rewriting it with -update.
func golden(t *testing.T, name string, got []byte) {
	t.Helper()
	if *update {
		require.NoError(t, os.MkdirAll(filepath.Dir(name), os.ModePerm))
		require.NoError(t, os.WriteFile(name, got, 0o644))
	}
	want, err := os.ReadFile(name)
	require.NoError(t, err, "run go test ./inject/ -update")
	assert.Equal(t, string(want), string(got), "%s differs; run go test ./inject/ -update and review the diff", name)
}

func TestGenerator_RenderImplementation(t *testing.T) {
	clock := at("stock/memory", "NewClock", nil, []string{"*Clock"}, provide("index=0"))
	order := consumer()
	order.Func.Parameters = append(order.Func.Parameters, annotation.EntryFuncType{Type: "*memory.Clock"})
	order.Imports["memory"] = testModule + "/stock/memory"
	entries := []annotation.Entry{impl("memory", provide("index=0"), deflt()), impl("unlimited", provide("index=0")), order, clock}

	files, err := NewGenerator(newTestGraph(t, entries), testOptions("/mod")).Render(context.Background())
	require.NoError(t, err)

	src := files["implementation/stock/service_implementation.go"]
	require.NotNil(t, src, "the switch file is rendered")
	golden(t, "testdata/implementation/service_implementation.go.golden", src)

	imp, call := Implementation{Type: "Service", Path: testModule + "/stock"}.Module("example.com/app/gen")
	service := string(files["order/newservice_module.go"])
	assert.Contains(t, service, imp.Alias+` "example.com/app/gen/implementation/stock"`)
	assert.Contains(t, service, `fx.Module("example.com/app/order.NewService",
			`+call+`,
			`+genAlias(clock)+`.NewClockModule(),
`, "the set goes through the switch, another value of a member's package keeps its module")
	assert.NotContains(t, service, ".NewModule()", "a consumer never wires a member directly")

	root := string(files["module.go"])
	assert.Equal(t, 1, strings.Count(root, call), "the root includes the switch once")
	assert.NotContains(t, root, ".NewModule()", "the root never wires a member directly")
	assert.Contains(t, root, genAlias(clock)+".NewClockModule(),")
	assert.Contains(t, root, imp.Alias+` "example.com/app/gen/implementation/stock"`)
}

func TestGenerator_RenderNamedImplementation(t *testing.T) {
	entries := []annotation.Entry{
		impl("memory", provide("name=formal,index=0"), deflt()),
		impl("unlimited", provide("name=formal,index=0")),
		consumer(inject("name=formal,index=0")),
	}

	files, err := NewGenerator(newTestGraph(t, entries), testOptions("/mod")).Render(context.Background())
	require.NoError(t, err)

	sw := string(files["implementation/stock/serviceformal_implementation.go"])
	assert.Contains(t, sw, `serviceFormalImplementationKey = "stock.formal.implementation"`)
	assert.Contains(t, sw, "func ServiceFormalImplementationModule() fx.Option {")

	_, call := Implementation{Type: "Service", Path: testModule + "/stock", Name: "formal"}.Module("example.com/app/gen")
	assert.Contains(t, string(files["order/newservice_module.go"]), call+",")
}
