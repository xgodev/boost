package otel

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	bootel "github.com/xgodev/boost/factory/contrib/go.opentelemetry.io/otel/v1"
	"github.com/xgodev/boost/factory/contrib/modelcontextprotocol/go-sdk/v1/client"
	"github.com/xgodev/boost/wrapper/config"
	koanf "github.com/xgodev/boost/wrapper/config/contrib/knadh/koanf/v1"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

// connect points the boost tracer provider at a recorder, loads config from
// env (boost's koanf keeps keys across Load calls, so each test sets
// .enabled) and connects a client with Register.
func connect(t *testing.T, enabled string) (*client.Client, *tracetest.SpanRecorder, error) {
	t.Helper()
	t.Setenv("BOOST_FACTORY_MCP_CLIENT_PLUGINS_OTEL_ENABLED", enabled)
	config.Set(koanf.New())
	config.Load()

	sr := tracetest.NewSpanRecorder()
	prev := bootel.TracerProvider
	bootel.TracerProvider = sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr))
	t.Cleanup(func() { bootel.TracerProvider = prev })

	srv := mcp.NewServer(&mcp.Implementation{Name: "upstream", Version: "1.0.0"}, nil)
	ts := httptest.NewServer(mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return srv }, nil))
	t.Cleanup(ts.Close)

	cli, err := client.NewClientWithOptions(context.Background(), &client.Options{
		Name: "test", Version: "1.0.0", Transport: client.TransportHTTP, Endpoint: ts.URL, MaxRetries: -1,
	}, Register)
	if err == nil {
		t.Cleanup(func() { cli.Close() })
	}
	return cli, sr, err
}

func clientSpans(sr *tracetest.SpanRecorder, name string) int {
	n := 0
	for _, s := range sr.Ended() {
		if s.Name() == name && s.SpanKind() == trace.SpanKindClient {
			n++
		}
	}
	return n
}

func TestRegister_TracesEachCall(t *testing.T) {
	cli, sr, err := connect(t, "true")
	require.NoError(t, err)

	_, err = cli.Session().ListTools(context.Background(), nil)
	require.NoError(t, err)

	assert.Equal(t, 1, clientSpans(sr, "initialize"), "the plugin runs before Connect")
	assert.Equal(t, 1, clientSpans(sr, "tools/list"))
	assert.Equal(t, Scope, sr.Ended()[0].InstrumentationScope().Name)
}

func TestRegister_Disabled(t *testing.T) {
	cli, sr, err := connect(t, "false")
	require.NoError(t, err)

	_, err = cli.Session().ListTools(context.Background(), nil)
	require.NoError(t, err)
	assert.Empty(t, sr.Ended())
}

func TestRegister_OptionsError(t *testing.T) {
	_, _, err := connect(t, "maybe")
	assert.Error(t, err)
}

func TestNewOptionsWithPath(t *testing.T) {
	t.Setenv("BOOST_FACTORY_MCP_CLIENT_PLUGINS_OTEL_ENABLED", "true")
	t.Setenv("BOOST_FACTORY_MCP_CLIENT_GITHUB_PLUGINS_OTEL_ENABLED", "false")
	config.Set(koanf.New())
	config.Load()

	o, err := NewOptionsWithPath("boost.factory.mcp.client.github.plugins.otel")
	require.NoError(t, err)
	assert.False(t, o.Enabled)
}
