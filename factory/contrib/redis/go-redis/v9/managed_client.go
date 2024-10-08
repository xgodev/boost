package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// ManagedClient represents a redis managed client.
type ManagedClient struct {
	Client  *redis.Client
	Plugins []Plugin
	Options *Options
}

// NewManagedClientWithConfigPath returns a new managed client with options from config path.
func NewManagedClientWithConfigPath(ctx context.Context, path string, plugins ...Plugin) (*ManagedClient, error) {

	opts, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}

	client, err := NewClientWithOptions(ctx, opts, plugins...)
	if err != nil {
		return nil, err
	}

	return &ManagedClient{
		Client:  client,
		Plugins: plugins,
		Options: opts,
	}, nil
}

// NewManagedClient returns a new managed client with default options.
func NewManagedClient(ctx context.Context, plugins ...Plugin) (*ManagedClient, error) {
	opts, err := NewOptions()
	if err != nil {
		return nil, err
	}

	client, err := NewClientWithOptions(ctx, opts, plugins...)
	if err != nil {
		return nil, err
	}

	return &ManagedClient{
		Client:  client,
		Plugins: plugins,
		Options: opts,
	}, nil
}

// NewManagedClientWithOptions returns a new managed client with options.
func NewManagedClientWithOptions(ctx context.Context, opts *Options, plugins ...Plugin) (*ManagedClient, error) {
	s, err := NewClientWithOptions(ctx, opts, plugins...)
	if err != nil {
		return nil, err
	}

	return &ManagedClient{
		Client:  s,
		Plugins: plugins,
		Options: opts,
	}, nil
}
