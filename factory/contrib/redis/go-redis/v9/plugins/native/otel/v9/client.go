package otel

import (
	"context"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/xgodev/boost/wrapper/log"
)

// Client represents a otel client for redis.
type Client struct {
	options *Options
}

// NewClientWithConfigPath returns a new otel client with options from config path.
func NewClientWithConfigPath(path string) (*Client, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}

	return NewClientWithOptions(o), nil
}

// NewClient returns a new otel client with default options.
func NewClient() (*Client, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}

	return NewClientWithOptions(o), nil
}

// NewClientWithOptions returns a new otel client with options.
func NewClientWithOptions(options *Options) *Client {
	return &Client{options: options}
}

// Register registers this otel client to redis client.
func (d *Client) Register(ctx context.Context, client *redis.Client) error {
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

// ClientRegister registers a new otel client to redis client.
func ClientRegister(ctx context.Context, client *redis.Client) error {
	o, err := NewOptions()
	if err != nil {
		return err
	}
	d := NewClientWithOptions(o)
	return d.Register(ctx, client)
}
