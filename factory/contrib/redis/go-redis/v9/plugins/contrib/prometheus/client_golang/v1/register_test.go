package prometheus

import (
	"context"
	"testing"
	"time"

	prom "github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// unreachable is an address nothing listens on, so every command fails fast
// and lands in the error counter when the hook is installed.
const unreachable = "127.0.0.1:1"

func failedGets(t *testing.T) float64 {
	t.Helper()
	m := &dto.Metric{}
	require.NoError(t, singleErrors.WithLabelValues("get").(prom.Metric).Write(m))
	return m.GetCounter().GetValue()
}

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

func TestClientRegisterInstrumentsCommands(t *testing.T) {
	ctx := context.Background()
	client := newClient(t)
	before := failedGets(t)

	require.NoError(t, NewClientWithOptions(&Options{Enabled: true}).Register(ctx, client))
	_ = client.Get(ctx, "k").Err()

	assert.Equal(t, before+1, failedGets(t))
}

func TestClientRegisterDisabledDoesNothing(t *testing.T) {
	ctx := context.Background()
	client := newClient(t)
	before := failedGets(t)

	require.NoError(t, NewClientWithOptions(&Options{Enabled: false}).Register(ctx, client))
	_ = client.Get(ctx, "k").Err()

	assert.Equal(t, before, failedGets(t))
}

func TestClusterRegisterInstrumentsCommands(t *testing.T) {
	ctx := context.Background()
	client := newCluster(t)
	before := failedGets(t)

	require.NoError(t, NewClusterWithOptions(&Options{Enabled: true}).Register(ctx, client))
	_ = client.Get(ctx, "k").Err()

	assert.Equal(t, before+1, failedGets(t))
}

func TestClusterRegisterDisabledDoesNothing(t *testing.T) {
	ctx := context.Background()
	client := newCluster(t)
	before := failedGets(t)

	require.NoError(t, NewClusterWithOptions(&Options{Enabled: false}).Register(ctx, client))
	_ = client.Get(ctx, "k").Err()

	assert.Equal(t, before, failedGets(t))
}
