package otel

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	otelboost "github.com/xgodev/boost/factory/contrib/go.opentelemetry.io/otel/v1"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// unreachable is an address nothing listens on: commands fail fast, but the
// tracing hook still records a span for each one.
const unreachable = "127.0.0.1:1"

func startTracing(t *testing.T) *tracetest.InMemoryExporter {
	t.Helper()
	boost.Start()
	ctx := context.Background()
	opts, err := otelboost.NewOptions()
	require.NoError(t, err)
	spans := tracetest.NewInMemoryExporter()
	otelboost.StartTracerProviderWithOptions(ctx, opts, sdktrace.WithSyncer(spans))
	return spans
}

func TestRegisterTracesClientAndClusterCommands(t *testing.T) {
	spans := startTracing(t)
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{Addr: unreachable, MaxRetries: -1, DialTimeout: 50 * time.Millisecond})
	t.Cleanup(func() { _ = client.Close() })
	cluster := redis.NewClusterClient(&redis.ClusterOptions{Addrs: []string{unreachable}, MaxRetries: -1, MaxRedirects: -1, DialTimeout: 50 * time.Millisecond})
	t.Cleanup(func() { _ = cluster.Close() })

	require.NoError(t, NewClientWithOptions(&Options{Enabled: true}).Register(ctx, client))
	_ = client.Get(ctx, "k").Err()
	assert.NotEmpty(t, spans.GetSpans(), "client command not traced")

	spans.Reset()
	require.NoError(t, NewClusterWithOptions(&Options{Enabled: true}).Register(ctx, cluster))
	_ = cluster.Get(ctx, "k").Err()
	assert.NotEmpty(t, spans.GetSpans(), "cluster command not traced")
}

func TestRegisterDisabledDoesNotTrace(t *testing.T) {
	spans := startTracing(t)
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{Addr: unreachable, MaxRetries: -1, DialTimeout: 50 * time.Millisecond})
	t.Cleanup(func() { _ = client.Close() })

	spans.Reset()
	require.NoError(t, NewClientWithOptions(&Options{Enabled: false}).Register(ctx, client))
	_ = client.Get(ctx, "k").Err()
	assert.Empty(t, spans.GetSpans())
}
