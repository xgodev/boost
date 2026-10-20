package redis

import (
	"testing"

	r "github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

// TestModules_Together: an app takes both the client and the cluster
// client; ValidateApp checks the graph without reaching redis.
func TestModules_Together(t *testing.T) {
	err := fx.ValidateApp(ClusterModule(), ClientModule(), fx.NopLogger,
		fx.Invoke(func(*r.Client, *r.ClusterClient) {}))
	if err != nil {
		t.Fatalf("ClientModule after ClusterModule: %v", err)
	}
}
