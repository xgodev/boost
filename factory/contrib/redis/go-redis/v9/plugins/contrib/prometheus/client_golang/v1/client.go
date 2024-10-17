package prometheus

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/xgodev/boost/wrapper/log"
)

// Client represents a datadog client for redis.
type Client struct {
	options *Options
}

// NewClientWithConfigPath returns a new datadog client with options from config path.
func NewClientWithConfigPath(path string) (*Client, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}

	return NewClientWithOptions(o), nil
}

// NewClient returns a new datadog client with default options.
func NewClient() (*Client, error) {
	o, err := NewOptions()
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

	client.AddHook(NewHook())

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
