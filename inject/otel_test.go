package inject

import (
	"bytes"
	"context"
	"go/parser"
	"go/token"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost/annotation"
)

// otelModulePath is the fx module that starts the OpenTelemetry providers;
// without it every OTel plugin reports to a noop provider.
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

// TestGenerator_RenderClientStartsOpenTelemetry renders a client of each
// factory that has an otel plugin, with and without it: the module of a
// client built with the otel plugin starts the providers.
func TestGenerator_RenderClientStartsOpenTelemetry(t *testing.T) {
	testCases := []struct {
		name, dir, pkg, path, tp, file string
		otel, plain                    string // the @Inject with and without the otel plugin
	}{
		{name: "resty", dir: "quote", pkg: "resty", path: restyPath, tp: "*resty.Client",
			otel: "index=0, name=quote, plugins=[log, otel]", plain: "index=0, name=quote, plugins=[log]",
			file: "client/github.com/go-resty/resty/v2/quote_client.go"},
		{name: "redis", dir: "redis", pkg: "redis", path: redisPath, tp: "*redis.Client",
			otel: "index=0, name=cache, plugins=[otel]", plain: "index=0, name=cache, plugins=[health]",
			file: "client/github.com/redis/go-redis/v9/cache_client.go"},
		{name: "mongo", dir: "mongo", pkg: "mongo", path: mongoPath, tp: "*mongo.Conn",
			otel: "index=0, name=stock, plugins=[otel]", plain: "index=0, name=stock, plugins=[health]",
			file: "client/github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v2/stock_conn.go"},
		{name: "sql", dir: "sql", pkg: "sql", path: "database/sql", tp: "*sql.DB",
			otel: "index=0, name=gold, drivers=[godror], plugins=[otel]", plain: "index=0, name=gold, drivers=[godror], plugins=[health]",
			file: "client/database/sql/gold_db.go"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for raw, want := range map[string]bool{tc.otel: true, tc.plain: false} {
				entry := at("client/"+tc.dir, "New", []string{tc.tp}, []string{"*Repository"}, provide("index=0"), inject(raw))
				entry.Imports = map[string]string{tc.pkg: tc.path}

				files, err := NewGenerator(newTestGraph(t, []annotation.Entry{entry}), testOptions("/mod")).Render(context.Background())
				require.NoError(t, err)

				src := files[tc.file]
				require.NotNil(t, src, "the client is rendered at %s", tc.file)
				assert.Equal(t, want, startsOpenTelemetry(t, tc.file, src), "@Inject(%s)", raw)
			}
		})
	}
}
