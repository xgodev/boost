package client

import (
	"context"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/xgodev/boost/wrapper/log"
	"os"
)

// NewServiceWithConfigPath returns a cache with options from config path .
func NewServiceWithConfigPath(ctx context.Context, path string) (dapr.Client, error) {
	options, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewServiceWithOptions(ctx, options)
}

// NewServiceWithOptions returns a cache with options.
func NewServiceWithOptions(ctx context.Context, o *Options) (dapr.Client, error) {

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

	return NewServiceWithOptions(ctx, o)
}
