package builder

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost/annotation"
	"github.com/xgodev/boost/model/errors"
)

// pkgDoc builds the entry of the package doc of package pkg of testModule.
func pkgDoc(pkg string, anns ...annotation.Annotation) annotation.Entry {
	return annotation.Entry{
		Path:        testModule + "/" + pkg,
		Package:     path.Base(pkg),
		File:        "doc.go",
		Annotations: anns,
		Imports:     fileImports,
	}
}

// ceFunction is a function app whose adapters are nats or cloudevents.
func ceFunction() []annotation.Entry {
	return []annotation.Entry{
		decl("events", "", "Handle", eventParams, eventResults, ann("Function", "adapter=[nats, cloudevents], app=worker")),
	}
}

func TestNewSpec_ServerDefaults(t *testing.T) {
	spec, err := newTestSpec(t, validEntries())
	require.NoError(t, err)

	api, worker := spec.Apps[0], spec.Apps[1]
	admin, grpcSrv, mcpSrv, ops, rest := api.Servers[0], api.Servers[1], api.Servers[2], api.Servers[3], api.Servers[4]
	restDefaults := []string{"recover", "requestId", "logger", "errorHandler", "health"}
	assert.Equal(t, restDefaults, rest.Builtin, "a rest server without @Server gets the default set")
	assert.Equal(t, restDefaults, admin.Builtin, "a named one too")
	assert.Equal(t, []string{"log"}, grpcSrv.Builtin)
	assert.Equal(t, []string{"log"}, mcpSrv.Builtin)
	assert.Equal(t, []string{"log"}, ops.Builtin)
	assert.Empty(t, worker.Function.Builtin, "the nats adapter has no plugins")

	spec, err = newTestSpec(t, ceFunction())
	require.NoError(t, err)
	assert.Empty(t, spec.Apps[0].Function.Builtin, "cloudevents has no default plugin")
}

func TestNewSpec_Server(t *testing.T) {
	entries := append(validEntries(),
		pkgDoc("server/rest",
			ann("Server", "type=rest, plugins=[recover, prometheus]"),
			ann("Server", "type=rest, name=admin, plugins=[]")),
		pkgDoc("server/grpc", ann("Server", "type=grpc, plugins=[prometheus, log]")),
		pkgDoc("server/mcp", ann("Server", "type=mcp, name=ops")),
	)
	spec, err := newTestSpec(t, append(entries, entries...)) // duplicates are collapsed
	require.NoError(t, err)

	api := spec.Apps[0]
	admin, grpcSrv, mcpSrv, ops, rest := api.Servers[0], api.Servers[1], api.Servers[2], api.Servers[3], api.Servers[4]
	assert.Equal(t, []string{"recover", "prometheus"}, rest.Builtin, "@Server lists the plugins in order")
	assert.Empty(t, admin.Builtin, "plugins=[] turns them all off")
	assert.Equal(t, []string{"prometheus", "log"}, grpcSrv.Builtin)
	assert.Equal(t, []string{"log"}, ops.Builtin, "@Server without plugins= keeps the default set")
	assert.Equal(t, []string{"log"}, mcpSrv.Builtin)
	assert.Empty(t, rest.Plugins, "no @Provide []Plugin")
}

func TestNewSpec_ServerOfAFunctionAdapter(t *testing.T) {
	spec, err := newTestSpec(t, append(ceFunction(), pkgDoc("server/events", ann("Server", "type=cloudevents, app=worker, plugins=[otel]"))))
	require.NoError(t, err)
	assert.Equal(t, map[string][]string{"cloudevents": {"otel"}}, spec.Apps[0].Function.Builtin)
}

func TestNewSpec_ServerProvidedPluginsStay(t *testing.T) {
	spec, err := newTestSpec(t, append(validEntries(), pluginEntries()...))
	require.NoError(t, err)
	for _, srv := range spec.Apps[0].Servers {
		assert.NotEmpty(t, srv.Plugins, srv.Name)
		assert.Empty(t, srv.Builtin, "server %s takes the []Plugin its @Provide gives, without the default set", srv.Name)
	}
}

func TestNewSpec_ServerErrors(t *testing.T) {
	for _, tc := range []struct {
		name    string
		entries []annotation.Entry
		is      func(error) bool
		want    string
	}{
		{"type is required", []annotation.Entry{pkgDoc("s", ann("Server", "plugins=[log]"))}, errors.IsNotValid, "type is required"},
		{"unknown type", []annotation.Entry{pkgDoc("s", ann("Server", "type=soap"))}, errors.IsNotValid, "cloudevents, grpc, mcp, rest"},
		{"unknown key", []annotation.Entry{pkgDoc("s", ann("Server", "type=rest, port=8080"))}, errors.IsNotValid, `unknown key "port"`},
		{"unknown plugin", []annotation.Entry{pkgDoc("s", ann("Server", "type=grpc, plugins=[recover]"))}, errors.IsNotValid, "plugin recover of grpc is not one of compressor, datadog, log, prometheus"},
		{"plugin twice", []annotation.Entry{pkgDoc("s", ann("Server", "type=rest, plugins=[logger, logger]"))}, errors.IsNotValid, "lists logger twice"},
		{"bad server name", []annotation.Entry{pkgDoc("s", ann("Server", "type=rest, name=Admin"))}, errors.IsNotValid, `server name "Admin"`},
		{"name of an adapter", []annotation.Entry{pkgDoc("s", ann("Server", "type=cloudevents, name=ce"))}, errors.IsNotValid, "name= does not apply"},
		{"on a function", []annotation.Entry{decl("s", "", "Plugins", nil, nil, ann("Server", "type=rest"))}, errors.IsNotValid, "package doc"},
		{"twice", []annotation.Entry{pkgDoc("a", ann("Server", "type=rest")), pkgDoc("b", ann("Server", "type=rest"))}, errors.IsAlreadyExists, "server rest of app api"},
		{"no such server", []annotation.Entry{pkgDoc("s", ann("Server", "type=rest, name=missing"))}, errors.IsNotFound, "app api has no server missing"},
		{"another type", []annotation.Entry{pkgDoc("s", ann("Server", "type=grpc, name=rest"))}, errors.IsNotValid, "server rest of app api serves rest"},
		{"adapter not used", []annotation.Entry{pkgDoc("s", ann("Server", "type=cloudevents, app=worker"))}, errors.IsNotFound, "app worker has no @Function with the cloudevents adapter"},
		{"with @Provide []Plugin", append(pluginEntries(), pkgDoc("s", ann("Server", "type=rest, name=admin"))), errors.IsAlreadyExists, "plugins.Admin"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := newTestSpec(t, append(validEntries(), tc.entries...))
			require.Error(t, err)
			assert.True(t, tc.is(err), "unexpected error kind: %v", err)
			assert.Contains(t, err.Error(), tc.want)
		})
	}
}

func TestNewSpec_OtherAnnotationsOnAPackageDocAreIgnored(t *testing.T) {
	_, err := newTestSpec(t, append(validEntries(), pkgDoc("s", ann("REST", "method=GET, path=/x"), ann("Middleware", ""))))
	require.NoError(t, err, "only @Server is read from a package doc")
}

// TestGenerator_RenderServer checks that each door builds its server with
// the plugins its @Server lists: the default server with the unnamed
// []Plugin of the shipped module, a named one with the value named after it
// and an adapter with the []Plugin of its fx module.
func TestGenerator_RenderServer(t *testing.T) {
	entries := append(validEntries(),
		pkgDoc("server/rest",
			ann("Server", "type=rest, plugins=[recover, prometheus]"),
			ann("Server", "type=rest, name=admin, plugins=[]")),
		pkgDoc("server/grpc", ann("Server", "type=grpc, plugins=[prometheus, log]")),
		pkgDoc("server/mcp", ann("Server", "type=mcp, name=ops, app=api, plugins=[log, otel]")),
		decl("events", "", "Handle", eventParams, eventResults, ann("Function", "adapter=cloudevents, app=events")),
		pkgDoc("server/events", ann("Server", "type=cloudevents, app=events, plugins=[otel]")),
	)
	files, err := render(t, entries)
	require.NoError(t, err)

	assert.NotContains(t, string(files["builder/api/admin_module.go"]), "plugins/", "plugins=[] builds the server without plugins")
	for _, name := range []string{"api/rest_module.go", "api/grpc_module.go", "api/mcp_module.go", "api/ops_module.go", "events/function_module.go"} {
		src, ok := files["builder/"+name]
		require.True(t, ok, "files: %v", fileNames(files))
		assertImportsUsed(t, name, src)

		golden := filepath.Join("testdata", "golden", "server", path.Base(name)+".golden")
		if path.Dir(name) == "events" {
			golden = filepath.Join("testdata", "golden", "server", "cloudevents_"+path.Base(name)+".golden")
		}
		if *update {
			require.NoError(t, os.MkdirAll(filepath.Dir(golden), 0o755))
			require.NoError(t, os.WriteFile(golden, src, 0o644))
			continue
		}
		want, err := os.ReadFile(golden)
		require.NoError(t, err, "run go test ./builder/ -run TestGenerator_RenderServer -update")
		assert.Equal(t, string(want), string(src), name)
	}
}
