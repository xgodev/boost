package client

import (
	"context"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/xgodev/boost/wrapper/log"
	"os"
)

// NewWithConfigPath returns a cache with options from config path .
func NewWithConfigPath(ctx context.Context, path string) (dapr.Client, error) {
	options, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewWithOptions(ctx, options)
}

// NewWithOptions returns a cache with options.
func NewWithOptions(ctx context.Context, o *Options) (dapr.Client, error) {

	logger := log.FromContext(ctx)

	s, err := dapr.NewClientWithAddressContext(ctx, o.Address)
	if err != nil {
		return nil, err
	}

	logger.Infof("Created dapr http service on address %v", o.Address)

	return s, nil
}

// NewService returns a cache.
func NewService(ctx context.Context) (dapr.Client, error) {

	o, err := NewOptions()
	if err != nil {
		return nil, err
	}

	if v := os.Getenv("DAPR_GRPC_ENDPOINT"); v != "" {
		o.Address = v
	}

	return NewWithOptions(ctx, o)
}
