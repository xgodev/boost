package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// ManagedClusterClient represents a redis cluster managed client.
type ManagedClusterClient struct {
	Client  *redis.ClusterClient
	Plugins []ClusterPlugin
	Options *Options
}

// NewManagedClusterClientWithConfigPath returns a new managed client with options from config path.
func NewManagedClusterClientWithConfigPath(ctx context.Context, path string, plugins ...ClusterPlugin) (*ManagedClusterClient, error) {

	opts, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}

	client, err := NewClusterClientWithOptions(ctx, opts, plugins...)
	if err != nil {
		return nil, err
	}

	return &ManagedClusterClient{
		Client:  client,
		Plugins: plugins,
		Options: opts,
	}, nil
}

// NewManagedClusterClient returns a new managed client with default options.
func NewManagedClusterClient(ctx context.Context, plugins ...ClusterPlugin) (*ManagedClusterClient, error) {
	opts, err := NewOptions()
	if err != nil {
		return nil, err
	}

	client, err := NewClusterClientWithOptions(ctx, opts, plugins...)
	if err != nil {
		return nil, err
	}

	return &ManagedClusterClient{
		Client:  client,
		Plugins: plugins,
		Options: opts,
	}, nil
}

// NewManagedClusterClientWithOptions returns a new managed client with options.
func NewManagedClusterClientWithOptions(ctx context.Context, opts *Options, plugins ...ClusterPlugin) (*ManagedClusterClient, error) {
	s, err := NewClusterClientWithOptions(ctx, opts, plugins...)
	if err != nil {
		return nil, err
	}

	return &ManagedClusterClient{
		Client:  s,
		Plugins: plugins,
		Options: opts,
	}, nil
}
