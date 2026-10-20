package builder

import (
	"bytes"
	"go/parser"
	"go/token"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost/annotation"
)

// otelModulePath is the fx module that starts the OpenTelemetry providers;
// without it every OTel plugin and middleware reports to a noop provider.
const otelModulePath = "github.com/xgodev/boost/fx/modules/factory/contrib/go.opentelemetry.io/otel/v1"

// startsOpenTelemetry tells whether src includes the module of otelModulePath.
func startsOpenTelemetry(t *testing.T, name string, src []byte) bool {
	t.Helper()
	file, err := parser.ParseFile(token.NewFileSet(), name, src, parser.ImportsOnly)
	require.NoError(t, err, "%s is not valid Go", name)
	for _, imp := range file.Imports {
		p, err := strconv.Unquote(imp.Path.Value)
		require.NoError(t, err)
		if p != otelModulePath {
			continue
		}
		alias := "otel"
		if imp.Name != nil {
			alias = imp.Name.Name
		}
		return bytes.Contains(src, []byte(alias+".Module()"))
	}
	return false
}

// TestGenerator_RenderStartsOpenTelemetry checks that each door using an
// OTel plugin or middleware includes the module that starts the providers,
// and that a door that uses none does not.
func TestGenerator_RenderStartsOpenTelemetry(t *testing.T) {
	testCases := []struct {
		name    string
		entries []annotation.Entry
		starts  []string // the files that start the providers
		plain   []string // the files that do not
	}{
		{
			name:    "rest server with the otel plugin",
			entries: append(validEntries(), pkgDoc("server/rest", ann("Server", "type=rest, plugins=[recover, otel]"))),
			starts:  []string{"builder/api/rest_module.go"},
			plain:   []string{"builder/api/admin_module.go", "builder/api/grpc_module.go", "builder/api/mcp_module.go"},
		},
		{
			name:    "mcp server with the otel plugin",
			entries: append(validEntries(), pkgDoc("server/mcp", ann("Server", "type=mcp, plugins=[log, otel]"))),
			starts:  []string{"builder/api/mcp_module.go"},
			plain:   []string{"builder/api/rest_module.go", "builder/api/admin_module.go", "builder/api/grpc_module.go"},
		},
		{
			name: "cloudevents adapter with the otel plugin",
			entries: append(validEntries(),
				decl("events", "", "Handle", eventParams, eventResults, ann("Function", "adapter=cloudevents, app=events")),
				pkgDoc("server/events", ann("Server", "type=cloudevents, app=events, plugins=[otel]"))),
			starts: []string{"builder/events/function_module.go"},
			plain:  []string{"builder/api/rest_module.go"},
		},
		{
			name: "function with the opentelemetry middleware",
			entries: append(validEntries(),
				decl("events", "", "Handle", eventParams, eventResults, ann("Function", "adapter=nats, app=events, middlewares=[opentelemetry]"))),
			starts: []string{"builder/events/function_module.go"},
			plain:  []string{"builder/api/rest_module.go", "builder/worker/function_module.go"},
		},
		{
			name: "nothing uses otel",
			entries: append(validEntries(),
				decl("events", "", "Handle", eventParams, eventResults, ann("Function", "adapter=nats, app=events, middlewares=[recovery]"))),
			plain: []string{"builder/api/rest_module.go", "builder/api/admin_module.go", "builder/api/grpc_module.go",
				"builder/api/mcp_module.go", "builder/worker/function_module.go", "builder/events/function_module.go"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			files, err := render(t, tc.entries)
			require.NoError(t, err)

			for name, src := range files {
				assertImportsUsed(t, name, src)
			}
			for _, name := range tc.starts {
				src, ok := files[name]
				require.True(t, ok, "files: %v", fileNames(files))
				assert.True(t, startsOpenTelemetry(t, name, src), "%s uses otel, so it starts the providers", name)
			}
			for _, name := range tc.plain {
				src, ok := files[name]
				require.True(t, ok, "files: %v", fileNames(files))
				assert.False(t, startsOpenTelemetry(t, name, src), "%s uses no otel", name)
			}
		})
	}
}
