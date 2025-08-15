package datadog

import (
	"context"

	redistrace "github.com/DataDog/dd-trace-go/contrib/redis/go-redis.v9/v2"
	"github.com/redis/go-redis/v9"
	datadog "github.com/xgodev/boost/factory/contrib/datadog/dd-trace-go/v1"
	"github.com/xgodev/boost/wrapper/log"
)

// Client represents a datadog client for redis.
type Client struct {
	options *Options
}

// NewClientWithConfigPath returns a new datadog client with options from config path.
func NewClientWithConfigPath(path string, traceOptions ...redistrace.ClientOption) (*Client, error) {
	o, err := NewOptionsWithPath(path, traceOptions...)
	if err != nil {
		return nil, err
	}

	if !datadog.IsTracerEnabled() {
		o.Enabled = false
	}

	return NewClientWithOptions(o), nil
}

// NewClient returns a new datadog client with default options.
func NewClient(traceOptions ...redistrace.ClientOption) (*Client, error) {
	o, err := NewOptions(traceOptions...)
	if err != nil {
		return nil, err
	}

	return NewClientWithOptions(o), nil
}

// NewClientWithOptions returns a new datadog client with options.
func NewClientWithOptions(options *Options) *Client {
	return &Client{options: options}
}

// Register registers this datadog client to redis client.
func (d *Client) Register(ctx context.Context, client *redis.Client) error {
	if !d.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("integrating redis in datadog")

	redistrace.WrapClient(client, d.options.TraceOptions...)

	logger.Debug("redis successfully integrated in datadog")

	return nil
}

// ClientRegister registers a new datadog client to redis client.
func ClientRegister(ctx context.Context, client *redis.Client) error {
	o, err := NewOptions()
	if err != nil {
		return err
	}
	d := NewClientWithOptions(o)
	return d.Register(ctx, client)
}
