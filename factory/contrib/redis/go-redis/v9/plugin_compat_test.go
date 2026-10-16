package redis

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/xgodev/boost/factory"
)

// Each plugin type must stay an alias of factory.Plugin: plugins written
// against either name keep compiling and stay interchangeable.
func TestPluginTypesAreFactoryPlugin(t *testing.T) {
	var _ factory.Plugin[*redis.Client] = Plugin(func(context.Context, *redis.Client) error { return nil })
	var _ factory.Plugin[*redis.ClusterClient] = ClusterPlugin(func(context.Context, *redis.ClusterClient) error { return nil })
}
