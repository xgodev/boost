package prometheus

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/xgodev/boost/wrapper/log"
)

// Cluster represents a datadog client for redis cluster client.
type Cluster struct {
	options *Options
}

// NewClusterWithConfigPath returns datadog client with options from config path.
func NewClusterWithConfigPath(path string) (*Cluster, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}

	return NewClusterWithOptions(o), nil
}

// NewCluster returns datadog client with default options.
func NewCluster() (*Cluster, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}

	return NewClusterWithOptions(o), nil
}

// NewClusterWithOptions returns datadog client with options.
func NewClusterWithOptions(options *Options) *Cluster {
	return &Cluster{options: options}
}

// Register registers this datadog client on redis cluster client.
func (d *Cluster) Register(ctx context.Context, client *redis.ClusterClient) error {
	if !d.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("integrating redis in datadog")

	client.AddHook(NewHook())

	logger.Debug("redis successfully integrated in datadog")

	return nil
}

// ClusterRegister registers a new datadog client on redis cluster client.
func ClusterRegister(ctx context.Context, client *redis.ClusterClient) error {
	o, err := NewOptions()
	if err != nil {
		return err
	}
	d := NewClusterWithOptions(o)
	return d.Register(ctx, client)
}
