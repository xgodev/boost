package datadog

import (
	"context"
	"testing"
	"time"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/mocktracer"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// unreachable is an address nothing listens on, so commands fail fast; a
// wrapped client still traces them.
const unreachable = "127.0.0.1:1"

func newClient(t *testing.T) *redis.Client {
	c := redis.NewClient(&redis.Options{Addr: unreachable, MaxRetries: -1, DialTimeout: 50 * time.Millisecond})
	t.Cleanup(func() { _ = c.Close() })
	return c
}

func newCluster(t *testing.T) *redis.ClusterClient {
	c := redis.NewClusterClient(&redis.ClusterOptions{Addrs: []string{unreachable}, MaxRetries: -1, MaxRedirects: -1, DialTimeout: 50 * time.Millisecond})
	t.Cleanup(func() { _ = c.Close() })
	return c
}

func spansAfterGet(t *testing.T, register func(context.Context) error, client redis.UniversalClient) int {
	t.Helper()
	mt := mocktracer.Start()
	defer mt.Stop()

	ctx := context.Background()
	require.NoError(t, register(ctx))
	_ = client.Get(ctx, "k").Err()

	return len(mt.FinishedSpans())
}

func TestClientRegisterTracesCommands(t *testing.T) {
	client := newClient(t)
	n := spansAfterGet(t, func(ctx context.Context) error {
		return NewClientWithOptions(&Options{Enabled: true}).Register(ctx, client)
	}, client)
	assert.Positive(t, n)
}

func TestClientRegisterDisabledDoesNothing(t *testing.T) {
	client := newClient(t)
	n := spansAfterGet(t, func(ctx context.Context) error {
		return NewClientWithOptions(&Options{Enabled: false}).Register(ctx, client)
	}, client)
	assert.Zero(t, n)
}

func TestClusterRegisterTracesCommands(t *testing.T) {
	client := newCluster(t)
	n := spansAfterGet(t, func(ctx context.Context) error {
		return NewClusterWithOptions(&Options{Enabled: true}).Register(ctx, client)
	}, client)
	assert.Positive(t, n)
}

func TestClusterRegisterDisabledDoesNothing(t *testing.T) {
	client := newCluster(t)
	n := spansAfterGet(t, func(ctx context.Context) error {
		return NewClusterWithOptions(&Options{Enabled: false}).Register(ctx, client)
	}, client)
	assert.Zero(t, n)
}
