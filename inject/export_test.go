package inject

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost/annotation"
	"github.com/xgodev/boost/model/errors"
)

func TestGeneratorOptions_BaseImport(t *testing.T) {
	base, err := testOptions("/mod").BaseImport()
	require.NoError(t, err)
	assert.Equal(t, "example.com/app/gen", base)

	nested, err := GeneratorOptions{ModulePath: testModule, ModuleDir: "/mod", OutputDir: "/mod/internal/gen"}.BaseImport()
	require.NoError(t, err)
	assert.Equal(t, "example.com/app/internal/gen", nested)

	_, err = GeneratorOptions{ModulePath: testModule, ModuleDir: "/mod", OutputDir: "/mod"}.BaseImport()
	assert.True(t, errors.IsNotValid(err), "unexpected error: %v", err)
}

func TestAlias(t *testing.T) {
	alias := Alias("pkg", "sdk-go", "github.com/cloudevents/sdk-go/v2")
	assert.True(t, strings.HasPrefix(alias, "pkg_sdk_go_"), alias)
	assert.Len(t, alias, len("pkg_sdk_go_")+8)
	assert.NotEqual(t, alias, Alias("pkg", "sdk-go", "github.com/cloudevents/sdk-go/v3"), "the path tells packages of the same name apart")
}

func TestGenPath(t *testing.T) {
	const base = testModule + "/gen"

	testCases := []struct {
		name, base, path, want string
	}{
		{name: "package of the project", path: testModule + "/greet", want: "greet"},
		{name: "leading internal stays", path: testModule + "/internal/app/api/rest", want: "internal/app/api/rest"},
		{name: "deeper internal is renamed", path: testModule + "/internal/pkg/x/internal/y", want: "internal/pkg/x/internal_/y"},
		{name: "project root", path: testModule, want: "app"},
		{name: "dir of the builder", path: testModule + "/builder/api", want: "builder_/api"},
		{name: "dir of the clients", path: testModule + "/client", want: "client_"},
		{name: "dir of the switches", path: testModule + "/implementation/stock", want: "implementation_/stock"},
		{name: "first element with a dot", path: testModule + "/github.com/x", want: "github.com_/x"},
		{name: "package of another module", path: "github.com/xgodev/boost/greet", want: "github.com/xgodev/boost/greet"},
		{name: "internal of another module", path: "github.com/xgodev/boost/internal/greet", want: "github.com/xgodev/boost/internal_/greet"},
		{name: "path sharing the prefix", path: testModule + "x/greet", want: testModule + "x/greet"},
		{name: "project below the module", base: testModule + "/internal/gen", path: testModule + "/internal/app/api", want: "app/api"},
		{name: "module package outside the project", base: testModule + "/internal/gen", path: testModule + "/pkg/x", want: testModule + "/pkg/x"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := tc.base
			if b == "" {
				b = base
			}
			assert.Equal(t, tc.want, GenPath(b, tc.path))
		})
	}
}

func TestGenPathUnder(t *testing.T) {
	const base = testModule + "/gen"

	assert.Equal(t, "implementation/stock", GenPathUnder(base, "implementation", testModule+"/stock"))
	assert.Equal(t, "implementation/internal_/pkg/x/internal_/y", GenPathUnder(base, "implementation", testModule+"/internal/pkg/x/internal/y"),
		"below implementation/ no internal sits right under the output dir")
	assert.Equal(t, "builder/middleware/github.com/xgodev/boost/internal_/greet", GenPathUnder(base, "builder/middleware", "github.com/xgodev/boost/internal/greet"))
}

func TestModuleOf(t *testing.T) {
	entry := at("greet", "NewGreeter", nil, []string{"*Greeter"}, provide("index=0"))

	imp, call := ModuleOf("example.com/app/gen", entry)

	assert.Equal(t, ImportData{Alias: genAlias(entry), Path: "example.com/app/gen/greet"}, imp)
	assert.Equal(t, genAlias(entry)+".NewGreeterModule()", call)
}

func TestProvider(t *testing.T) {
	g := newTestGraph(t, append(greetingEntries(),
		at("probe", "NewProbe", nil, []string{"Probe", "error"}, provide("index=0"))))

	testCases := []struct {
		name     string
		pkg, tp  string
		wantKey  string
		wantIdx  int
		wantFind bool
	}{
		{name: "local type", pkg: "greet", tp: "*Greeter", wantKey: "example.com/app/greet.NewGreeter", wantFind: true},
		{name: "qualified from another package", pkg: "app", tp: "*greet.Greeter", wantKey: "example.com/app/greet.NewGreeter", wantFind: true},
		{name: "value type", pkg: "probe", tp: "Probe", wantKey: "example.com/app/probe.NewProbe", wantFind: true},
		{name: "pointer is another type", pkg: "probe", tp: "*Probe"},
		{name: "grouped values are not default providers", pkg: "names", tp: "string"},
		{name: "unknown type", pkg: "app", tp: "*Missing"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			comp, idx, ok := Provider(g, tc.pkg, tc.tp)
			require.Equal(t, tc.wantFind, ok)
			if ok {
				assert.Equal(t, tc.wantKey, comp.Key())
				assert.Equal(t, tc.wantIdx, idx)
			}
		})
	}
}

func TestProviderFor(t *testing.T) {
	received := at("pricereceived/boost", "NewPublisher", nil, []string{"*Publisher"}, provide("index=0"))
	changed := at("pricechanged/boost", "NewPublisher", nil, []string{"*Publisher"}, provide("index=0"))
	g := newTestGraph(t, []annotation.Entry{received, changed})

	consumer := at("app", "Run", nil, nil)
	consumer.Imports = map[string]string{"boost": testModule + "/pricechanged/boost"}

	testCases := []struct {
		name    string
		entry   annotation.Entry
		tp      string
		wantKey string
	}{
		{name: "local type in one of two same-named packages", entry: received, tp: "*Publisher",
			wantKey: testModule + "/pricereceived/boost.NewPublisher"},
		{name: "the other one", entry: changed, tp: "*Publisher",
			wantKey: testModule + "/pricechanged/boost.NewPublisher"},
		{name: "qualified, resolved by the file's import", entry: consumer, tp: "*boost.Publisher",
			wantKey: testModule + "/pricechanged/boost.NewPublisher"},
		{name: "unknown type", entry: consumer, tp: "*Missing"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			comp, idx, ok := ProviderFor(g, tc.entry, tc.tp)
			require.Equal(t, tc.wantKey != "", ok)
			if ok {
				assert.Equal(t, tc.wantKey, comp.Key())
				assert.Equal(t, 0, idx)
			}
		})
	}
}

func TestExecuteTemplate(t *testing.T) {
	tmpl := template.Must(template.New("t").Parse(`{{define "ok"}}package {{.}}
func  A( ) {}{{end}}{{define "bad"}}package {{.}} func{{end}}{{define "fail"}}{{.Missing}}{{end}}`))

	src, err := ExecuteTemplate(tmpl, "ok", "p")
	require.NoError(t, err)
	assert.Equal(t, "package p\n\nfunc A() {}\n", string(src), "the output is gofmt'd")

	_, err = ExecuteTemplate(tmpl, "bad", "p")
	assert.Error(t, err, "invalid Go is reported")

	_, err = ExecuteTemplate(tmpl, "fail", "p")
	assert.Error(t, err, "template errors are reported")
}

func TestWriteGenerated(t *testing.T) {
	const header = "// Code generated by other; DO NOT EDIT."
	dir := t.TempDir()
	write := func(rel, content string) {
		p := filepath.Join(dir, filepath.FromSlash(rel))
		require.NoError(t, os.MkdirAll(filepath.Dir(p), os.ModePerm))
		require.NoError(t, os.WriteFile(p, []byte(content), 0o644))
	}
	write("old/stale.go", header+"\n\npackage old\n")
	write("inject/kept.go", GeneratedHeader+"\n\npackage inject\n")

	files := map[string][]byte{"api/modules.go": []byte(header + "\n\npackage api\n")}
	require.NoError(t, WriteGenerated(context.Background(), dir, files, header))

	got, err := os.ReadFile(filepath.Join(dir, "api", "modules.go"))
	require.NoError(t, err)
	assert.Equal(t, string(files["api/modules.go"]), string(got))
	assert.NoDirExists(t, filepath.Join(dir, "old"), "stale files with the same header are removed")
	assert.FileExists(t, filepath.Join(dir, "inject", "kept.go"), "files of another generator are kept")
}
