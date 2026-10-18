package otel

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	bootel "github.com/xgodev/boost/factory/contrib/go.opentelemetry.io/otel/v1"
	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/server"
	"github.com/xgodev/boost/wrapper/config"
	koanf "github.com/xgodev/boost/wrapper/config/contrib/knadh/koanf/v1"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

// setup points the boost tracer provider at a recorder and loads config from
// env; boost's koanf keeps keys across Load calls, so each test sets .enabled.
func setup(t *testing.T, enabled string) (*server.Server, *tracetest.SpanRecorder) {
	t.Helper()
	t.Setenv("BOOST_FACTORY_MCP_SERVER_PLUGINS_OTEL_ENABLED", enabled)
	config.Set(koanf.New())
	config.Load()

	sr := tracetest.NewSpanRecorder()
	prev := bootel.TracerProvider
	bootel.TracerProvider = sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr))
	t.Cleanup(func() { bootel.TracerProvider = prev })

	srv, err := server.NewServerWithOptions(context.Background(),
		&server.Options{Name: "test", Version: "1.0.0", Transport: server.TransportHTTP})
	require.NoError(t, err)
	return srv, sr
}

func listTools(t *testing.T, srv *server.Server) {
	t.Helper()
	ctx := context.Background()
	st, ct := mcp.NewInMemoryTransports()
	ss, err := srv.Instance().Connect(ctx, st, nil)
	require.NoError(t, err)
	defer ss.Close()

	cs, err := mcp.NewClient(&mcp.Implementation{Name: "c", Version: "1.0.0"}, nil).Connect(ctx, ct, nil)
	require.NoError(t, err)
	defer cs.Close()
	_, err = cs.ListTools(ctx, nil)
	require.NoError(t, err)
}

func serverSpans(sr *tracetest.SpanRecorder, name string) int {
	n := 0
	for _, s := range sr.Ended() {
		if s.Name() == name && s.SpanKind() == trace.SpanKindServer {
			n++
		}
	}
	return n
}

func TestRegister_TracesEachCall(t *testing.T) {
	srv, sr := setup(t, "true")

	require.NoError(t, Register(context.Background(), srv))
	listTools(t, srv)

	assert.Equal(t, 1, serverSpans(sr, "tools/list"))
	assert.Equal(t, Scope, sr.Ended()[0].InstrumentationScope().Name)
}

func TestRegister_Disabled(t *testing.T) {
	srv, sr := setup(t, "false")

	require.NoError(t, Register(context.Background(), srv))
	listTools(t, srv)

	assert.Empty(t, sr.Ended())
}

func TestRegister_OptionsError(t *testing.T) {
	srv, _ := setup(t, "maybe")

	assert.Error(t, Register(context.Background(), srv))
}

func TestNewOptionsWithPath(t *testing.T) {
	t.Setenv("BOOST_FACTORY_MCP_SERVER_PLUGINS_OTEL_ENABLED", "true")
	t.Setenv("BOOST_FACTORY_MCP_SERVER_ADMIN_PLUGINS_OTEL_ENABLED", "false")
	config.Set(koanf.New())
	config.Load()

	o, err := NewOptionsWithPath("boost.factory.mcp.server.admin.plugins.otel")
	require.NoError(t, err)
	assert.False(t, o.Enabled)
}
