package redis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	vault "github.com/xgodev/boost/factory/contrib/mittwald/vaultgo/v0"
	"github.com/xgodev/boost/factory/contrib/redis/go-redis/v9"
)

func TestConfigureNeedsThePassword(t *testing.T) {
	o := &redis.Options{Password: "old"}
	for _, m := range []interface {
		Configure(context.Context, map[string]interface{}) error
	}{
		NewManagerWithOptions(&vault.ManagerOptions{}, o),
		NewClusterManagerWithOptions(&vault.ManagerOptions{}, o),
	} {
		assert.Error(t, m.Configure(context.Background(), map[string]interface{}{"username": "u"}))
	}
	assert.Equal(t, "old", o.Password)
}
