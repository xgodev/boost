package datadog

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

const unreachable = "127.0.0.1:1"

func TestRegisterAcceptsClientAndCluster(t *testing.T) {
	ctx := context.Background()
	for _, enabled := range []bool{true, false} {
		client := redis.NewClient(&redis.Options{Addr: unreachable})
		cluster := redis.NewClusterClient(&redis.ClusterOptions{Addrs: []string{unreachable}})

		assert.NoError(t, NewClientWithOptions(&Options{Enabled: enabled}).Register(ctx, client))
		assert.NoError(t, NewClusterWithOptions(&Options{Enabled: enabled}).Register(ctx, cluster))

		_ = client.Close()
		_ = cluster.Close()
	}
}
