package redis

import (
	"context"

	goredis "github.com/redis/go-redis/v9"
	vault "github.com/xgodev/boost/factory/contrib/mittwald/vaultgo/v0"
	"github.com/xgodev/boost/factory/contrib/redis/go-redis/v9"
)

const root = vault.ManagersRoot + ".redis"

func init() {
	vault.ManagerConfigAdd(root)
}

// Manager rebuilds a redis client with the password vault hands it.
type Manager = vault.ClientManager[*goredis.Client]

// ClusterManager rebuilds a redis cluster client with the password vault hands it.
type ClusterManager = vault.ClientManager[*goredis.ClusterClient]

// NewManager returns a redis manager with default manager options.
func NewManager(o *redis.Options, plugins ...redis.Plugin) (*Manager, error) {
	return NewManagerWithConfigPath(root, o, plugins...)
}

// NewManagerWithConfigPath returns a redis manager with manager options from config path.
func NewManagerWithConfigPath(path string, o *redis.Options, plugins ...redis.Plugin) (*Manager, error) {
	mo, err := vault.NewManagerOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewManagerWithOptions(mo, o, plugins...), nil
}

// NewManagerWithOptions returns a redis manager. Each Configure builds a
// client from a copy of o with the vault password; o is not changed.
func NewManagerWithOptions(mo *vault.ManagerOptions, o *redis.Options, plugins ...redis.Plugin) *Manager {
	return vault.NewClientManager(mo,
		func(ctx context.Context, data map[string]interface{}) (*goredis.Client, error) {
			password, err := vault.String(data, "password")
			if err != nil {
				return nil, err
			}
			opts := *o
			opts.Password = password
			return redis.NewClientWithOptions(ctx, &opts, plugins...)
		},
		func(_ context.Context, c *goredis.Client) error { return c.Close() },
	)
}

// NewClusterManager returns a redis cluster manager with default manager options.
func NewClusterManager(o *redis.Options, plugins ...redis.ClusterPlugin) (*ClusterManager, error) {
	return NewClusterManagerWithConfigPath(root, o, plugins...)
}

// NewClusterManagerWithConfigPath returns a redis cluster manager with manager options from config path.
func NewClusterManagerWithConfigPath(path string, o *redis.Options, plugins ...redis.ClusterPlugin) (*ClusterManager, error) {
	mo, err := vault.NewManagerOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewClusterManagerWithOptions(mo, o, plugins...), nil
}

// NewClusterManagerWithOptions returns a redis cluster manager. Each
// Configure builds a client from a copy of o with the vault password.
func NewClusterManagerWithOptions(mo *vault.ManagerOptions, o *redis.Options, plugins ...redis.ClusterPlugin) *ClusterManager {
	return vault.NewClientManager(mo,
		func(ctx context.Context, data map[string]interface{}) (*goredis.ClusterClient, error) {
			password, err := vault.String(data, "password")
			if err != nil {
				return nil, err
			}
			opts := *o
			opts.Password = password
			return redis.NewClusterClientWithOptions(ctx, &opts, plugins...)
		},
		func(_ context.Context, c *goredis.ClusterClient) error { return c.Close() },
	)
}
