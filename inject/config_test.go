package inject

import (
	"context"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost/annotation"
	"github.com/xgodev/boost/model/errors"
)

func configAnn(value string) annotation.Annotation { return annotation.NewAnnotation("Config", value) }

// cfg builds the entry of a @Config struct declared in package pkg of
// testModule, whose file imports time.
func cfg(pkg, name string, fields ...annotation.EntryField) annotation.Entry {
	return annotation.Entry{
		Path:        testModule + "/" + pkg,
		Package:     filepath.Base(pkg),
		Struct:      name,
		Fields:      fields,
		Annotations: []annotation.Annotation{configAnn("")},
		Imports:     map[string]string{"time": "time"},
	}
}

func field(name, tp, tag string) annotation.EntryField {
	return annotation.EntryField{Name: name, Type: tp, Tag: tag}
}

func TestNewComponent_Config(t *testing.T) {
	testCases := []struct {
		name  string
		field annotation.EntryField
		want  ConfigField
	}{
		{name: "string with an acronym", field: field("BaseURL", "string", `default:"http://x" desc:"d"`),
			want: ConfigField{Name: "BaseURL", Segment: "baseUrl", Type: "string", Default: `"http://x"`, Desc: "d"}},
		{name: "leading acronym", field: field("HTTPTimeout", "time.Duration", `default:"2s" desc:"d"`),
			want: ConfigField{Name: "HTTPTimeout", Segment: "httpTimeout", Type: "time.Duration", Default: "2 * time.Second", Desc: "d"}},
		{name: "duration in the largest whole unit", field: field("Wait", "time.Duration", `default:"1m30s" desc:"d"`),
			want: ConfigField{Name: "Wait", Segment: "wait", Type: "time.Duration", Default: "90 * time.Second", Desc: "d"}},
		{name: "duration zero value", field: field("Wait", "time.Duration", `desc:"d"`),
			want: ConfigField{Name: "Wait", Segment: "wait", Type: "time.Duration", Default: "time.Duration(0)", Desc: "d"}},
		{name: "int8 keeps its type", field: field("N", "int8", `default:"5" desc:"d"`),
			want: ConfigField{Name: "N", Segment: "n", Type: "int8", Default: "int8(5)", Desc: "d"}},
		{name: "int", field: field("MaxAttempts", "int", `default:"7" desc:"d"`),
			want: ConfigField{Name: "MaxAttempts", Segment: "maxAttempts", Type: "int", Default: "7", Desc: "d"}},
		{name: "uint zero value", field: field("N", "uint16", `desc:"d"`),
			want: ConfigField{Name: "N", Segment: "n", Type: "uint16", Default: "uint16(0)", Desc: "d"}},
		{name: "float", field: field("Ratio", "float64", `default:"1.5" desc:"d"`),
			want: ConfigField{Name: "Ratio", Segment: "ratio", Type: "float64", Default: "float64(1.5)", Desc: "d"}},
		{name: "bool", field: field("On", "bool", `default:"true" desc:"d"`),
			want: ConfigField{Name: "On", Segment: "on", Type: "bool", Default: "true", Desc: "d"}},
		{name: "string list", field: field("Items", "[]string", `default:"a, b" desc:"d"`),
			want: ConfigField{Name: "Items", Segment: "items", Type: "[]string", Default: `[]string{"a", "b"}`, Desc: "d"}},
		{name: "empty string list", field: field("Items", "[]string", `default:"" desc:"d"`),
			want: ConfigField{Name: "Items", Segment: "items", Type: "[]string", Default: "[]string{}", Desc: "d"}},
		{name: "string zero value", field: field("Name", "string", `desc:"d"`),
			want: ConfigField{Name: "Name", Segment: "name", Type: "string", Default: `""`, Desc: "d"}},
		{name: "config tag names the segment", field: field("Name", "string", `config:"other" desc:"d"`),
			want: ConfigField{Name: "Name", Segment: "other", Type: "string", Default: `""`, Desc: "d"}},
		{name: "hidden", field: field("Token", "string", `hide:"true" desc:"d"`),
			want: ConfigField{Name: "Token", Segment: "token", Type: "string", Default: `""`, Desc: "d", Hide: true}},
		{name: "underscore splits words", field: field("Max_Attempts", "int", `desc:"d"`),
			want: ConfigField{Name: "Max_Attempts", Segment: "maxAttempts", Type: "int", Default: "0", Desc: "d"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			comp, ok, err := newComponent(cfg("quote", "Options", tc.field))
			require.NoError(t, err)
			require.True(t, ok)
			assert.True(t, comp.IsConfig())
			assert.Equal(t, AnnotationTypePROVIDE, comp.Type)
			assert.Equal(t, "*Options", comp.Result(0))
			assert.Equal(t, testModule+"/quote.Options", comp.Key())
			assert.Equal(t, []ConfigField{tc.want}, comp.Config)
		})
	}
}

func TestNewComponent_ConfigSkipsAndAliases(t *testing.T) {
	entry := cfg("quote", "Options",
		field("Skipped", "map[string]string", `config:"-"`),
		field("internal", "int", ""),
		field("Wait", "t.Duration", `desc:"d"`),
	)
	entry.Imports = map[string]string{"t": "time"}

	comp, ok, err := newComponent(entry)
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, []ConfigField{
		{Name: "Wait", Segment: "wait", Type: "time.Duration", Default: "time.Duration(0)", Desc: "d"},
	}, comp.Config, `config:"-" and untagged unexported fields are not keys; time through an alias is a duration`)
}

func TestNewComponent_ConfigErrors(t *testing.T) {
	withAnn := func(entry annotation.Entry, anns ...annotation.Annotation) annotation.Entry {
		entry.Annotations = anns
		return entry
	}
	inMain := cfg("cmd", "Options")
	inMain.Package = "main"

	testCases := []struct {
		name  string
		entry annotation.Entry
		want  string
	}{
		{name: "missing desc", entry: cfg("quote", "Options", field("URL", "string", `default:"x"`)), want: "URL"},
		{name: "int out of range", entry: cfg("quote", "Options", field("N", "int8", `default:"300" desc:"d"`)), want: "N"},
		{name: "bad bool", entry: cfg("quote", "Options", field("On", "bool", `default:"yes" desc:"d"`)), want: "On"},
		{name: "bad duration", entry: cfg("quote", "Options", field("Wait", "time.Duration", `default:"2" desc:"d"`)), want: "Wait"},
		{name: "bad hide", entry: cfg("quote", "Options", field("Token", "string", `hide:"yes" desc:"d"`)), want: "Token"},
		{name: "unsupported type", entry: cfg("quote", "Options", field("Headers", "map[string]string", `desc:"d"`)), want: "map[string]string"},
		{name: "named type", entry: cfg("quote", "Options", field("Mode", "Mode", `desc:"d"`)), want: "Mode"},
		{name: "duration of another package", entry: cfg("quote", "Options", field("Wait", "other.Duration", `desc:"d"`)), want: "other.Duration"},
		{name: "embedded field", entry: cfg("quote", "Options", field("", "Base", "")), want: "embedded"},
		{name: "tagged unexported field", entry: cfg("quote", "Options", field("url", "string", `desc:"d"`)), want: "url"},
		{name: "segment not env reachable", entry: cfg("quote", "Options", field("URL", "string", `config:"base_url" desc:"d"`)), want: "base_url"},
		{name: "arguments", entry: withAnn(cfg("quote", "Options"), configAnn("root=x")), want: "no arguments"},
		{name: "on a function", entry: withAnn(at("quote", "New", nil, []string{"*Options"}), configAnn("")), want: "struct"},
		{name: "with another inject annotation", entry: withAnn(cfg("quote", "Options"), configAnn(""), provide("index=0")), want: "@Provide"},
		{name: "unexported struct", entry: cfg("quote", "options"), want: "exported"},
		{name: "package main", entry: inMain, want: "main"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := newComponent(tc.entry)
			require.Error(t, err)
			assert.True(t, errors.IsNotValid(err), "unexpected error: %v", err)
			assert.Contains(t, err.Error(), tc.want)
		})
	}
}

func TestNewGraph_Config(t *testing.T) {
	options := cfg("quote", "Options", field("URL", "string", `desc:"d"`))
	service := at("quote", "New", []string{"*Options"}, []string{"*Service"}, provide("index=0"))
	client := at("client", "New", []string{"*quote.Options"}, []string{"*Client"}, provide("index=0"))

	g := newTestGraph(t, []annotation.Entry{options, service, client})

	assert.Equal(t, []string{testModule + "/client.New", testModule + "/quote.New"},
		edgesOf(g.Edges[testModule+"/quote.Options"]),
		"constructors receiving *Options, in its package or qualified, depend on the config")

	t.Run("one @Config per package", func(t *testing.T) {
		_, err := NewGraphFromEntries(context.Background(), []annotation.Entry{options, cfg("quote", "Other")})
		require.Error(t, err)
		assert.True(t, errors.IsNotValid(err), "unexpected error: %v", err)
		assert.Contains(t, err.Error(), testModule+"/quote")
	})

	t.Run("a function providing the same struct", func(t *testing.T) {
		_, err := NewGraphFromEntries(context.Background(), []annotation.Entry{options,
			at("quote", "NewOptions", nil, []string{"*Options"}, provide("index=0"))})
		assert.True(t, errors.IsAlreadyExists(err), "unexpected error: %v", err)
	})

	t.Run("other inject annotations on a struct are still ignored", func(t *testing.T) {
		plain := cfg("quote", "Plain")
		plain.Annotations = []annotation.Annotation{provide("index=0")}
		g, err := NewGraphFromEntries(context.Background(), []annotation.Entry{plain})
		require.NoError(t, err)
		assert.Empty(t, g.Vertices)
	})
}

func TestGenerator_RenderConfig(t *testing.T) {
	options := cfg("quote", "Options",
		field("BaseURL", "string", `default:"http://x" desc:"base URL"`),
		field("Timeout", "time.Duration", `default:"2s" desc:"call timeout"`),
		field("Token", "string", `hide:"true" desc:"API token"`),
	)
	service := at("quote", "New", []string{"*Options"}, []string{"*Service"}, provide("index=0"))
	client := at("client", "New", []string{"*quote.Options"}, []string{"*Client"}, provide("index=0"))
	entries := []annotation.Entry{options, service, client}

	files, err := NewGenerator(newTestGraph(t, entries), testOptions("/mod")).Render(context.Background())
	require.NoError(t, err)

	src := string(files["quote/options_module.go"])
	require.NotEmpty(t, src, "the config module is generated next to the package's other modules")
	assert.True(t, strings.HasPrefix(src, generatedHeader))
	_, err = parser.ParseFile(token.NewFileSet(), "options_module.go", src, parser.AllErrors)
	require.NoError(t, err)

	for _, want := range []string{
		"package quote\n",
		"\t\"time\"\n",
		`const optionsRoot = "quote"`,
		"// env QUOTE_BASE__URL\n",
		`config.Add(optionsRoot+".baseUrl", "http://x", "base URL")`,
		"// env QUOTE_TIMEOUT\n",
		`config.Add(optionsRoot+".timeout", 2*time.Second, "call timeout")`,
		`config.Add(optionsRoot+".token", "", "API token", config.WithHide())`,
		"func OptionsModule() fx.Option {",
		"fx.Provide(func() (*" + srcAlias(options) + ".Options, error) {",
		"return config.NewOptionsWithPath[" + srcAlias(options) + ".Options](optionsRoot)",
	} {
		assert.Contains(t, src, want)
	}

	assert.Contains(t, string(files["quote/new_module.go"]), "\t\t\tOptionsModule(),\n",
		"a consumer in the package includes the config module")
	assert.Contains(t, string(files["client_/new_module.go"]), genAlias(options)+".OptionsModule(),",
		"a consumer in another package imports it")
	assert.Contains(t, string(files["module.go"]), genAlias(options)+".OptionsModule(),")

	t.Run("no duration, no time import", func(t *testing.T) {
		files, err := NewGenerator(newTestGraph(t, []annotation.Entry{cfg("quote", "Options")}), testOptions("/mod")).Render(context.Background())
		require.NoError(t, err)
		assert.NotContains(t, string(files["quote/options_module.go"]), `"time"`)
	})
}

func TestGenerator_RenderConfigKeyCollisions(t *testing.T) {
	testCases := []struct {
		name    string
		entries []annotation.Entry
		want    string
	}{
		{
			name: "a config key is an implementation key",
			entries: []annotation.Entry{
				impl("memory", provide("index=0"), deflt()), impl("unlimited", provide("index=0")),
				cfg("stock", "Options", field("Implementation", "string", `desc:"d"`)),
			},
			want: "stock.implementation",
		},
		{
			name: "a config key is the parent of another",
			entries: []annotation.Entry{
				cfg("stock", "Options", field("Memory", "string", `desc:"d"`)),
				cfg("stock/memory", "Options", field("Size", "int", `desc:"d"`)),
			},
			want: "stock.memory",
		},
		{
			name: "two packages map to the same root",
			entries: []annotation.Entry{
				cfg("stock_memory", "Options", field("Size", "int", `desc:"d"`)),
				cfg("stock/memory", "Options", field("Size", "int", `desc:"d"`)),
			},
			want: "stock.memory.size",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewGenerator(newTestGraph(t, tc.entries), testOptions("/mod")).Render(context.Background())
			require.Error(t, err)
			assert.True(t, errors.IsAlreadyExists(err), "unexpected error: %v", err)
			assert.Contains(t, err.Error(), tc.want)
		})
	}
}

func TestConfigKeyAndEnv(t *testing.T) {
	assert.Equal(t, "internal.pkg.client.quote.resty", ConfigKey("example.com/app", "example.com/app/internal/pkg/client/quote/resty"))
	assert.Equal(t, "app", ConfigKey("example.com/app", "example.com/app"))
	assert.Equal(t, "github.com.other.sdkGo.client", ConfigKey("example.com/app", "github.com/other/sdk-go/client"))

	assert.Equal(t, "INTERNAL_PKG_CLIENT_QUOTE_RESTY_BASE__URL", ConfigEnv("internal.pkg.client.quote.resty", "baseUrl"))
	assert.Equal(t, "STOCK_SDK__GO_TIMEOUT", ConfigEnv("stock.sdkGo", "timeout"))
}
