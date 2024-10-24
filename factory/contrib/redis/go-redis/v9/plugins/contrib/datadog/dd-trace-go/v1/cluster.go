package datadog

import (
	"context"

	datadog "github.com/xgodev/boost/factory/contrib/datadog/dd-trace-go/v1"

	"github.com/redis/go-redis/v9"
	"github.com/xgodev/boost/wrapper/log"
	redistrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/redis/go-redis.v9"
)

// Cluster represents a datadog client for redis cluster client.
type Cluster struct {
	options *Options
}

// NewClusterWithConfigPath returns datadog client with options from config path.
func NewClusterWithConfigPath(path string, traceOptions ...redistrace.ClientOption) (*Cluster, error) {
	o, err := NewOptionsWithPath(path, traceOptions...)
	if err != nil {
		return nil, err
	}

	if !datadog.IsTracerEnabled() {
		o.Enabled = false
	}

	return NewClusterWithOptions(o), nil
}

// NewCluster returns datadog client with default options.
func NewCluster(traceOptions ...redistrace.ClientOption) (*Cluster, error) {
	o, err := NewOptions(traceOptions...)
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

	redistrace.WrapClient(client, d.options.TraceOptions...)

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
