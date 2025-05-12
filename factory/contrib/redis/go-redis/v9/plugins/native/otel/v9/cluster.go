package otel

import (
	"context"
	"github.com/redis/go-redis/extra/redisotel/v9"

	"github.com/redis/go-redis/v9"
	"github.com/xgodev/boost/wrapper/log"
)

// Cluster represents a otel client for redis cluster client.
type Cluster struct {
	options *Options
}

// NewClusterWithConfigPath returns otel client with options from config path.
func NewClusterWithConfigPath(path string) (*Cluster, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}

	return NewClusterWithOptions(o), nil
}

// NewCluster returns otel client with default options.
func NewCluster() (*Cluster, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}

	return NewClusterWithOptions(o), nil
}

// NewClusterWithOptions returns otel client with options.
func NewClusterWithOptions(options *Options) *Cluster {
	return &Cluster{options: options}
}

// Register registers this otel client on redis cluster client.
func (d *Cluster) Register(ctx context.Context, client *redis.ClusterClient) error {
	if !d.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("integrating redis in otel")

	if err := redisotel.InstrumentMetrics(client); err != nil {
		return err
	}

	logger.Debug("redis successfully integrated in otel")

	return nil
}

// ClusterRegister registers a new otel client on redis cluster client.
func ClusterRegister(ctx context.Context, client *redis.ClusterClient) error {
	o, err := NewOptions()
	if err != nil {
		return err
	}
	d := NewClusterWithOptions(o)
	return d.Register(ctx, client)
}
