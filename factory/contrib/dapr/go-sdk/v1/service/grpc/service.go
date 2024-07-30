package grpc

import (
	"context"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/grpc"
	"github.com/xgodev/boost/wrapper/log"
)

// NewServiceWithConfigPath returns a cache with options from config path .
func NewServiceWithConfigPath(ctx context.Context, path string) (common.Service, error) {
	options, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewServiceWithOptions(ctx, options)
}

// NewServiceWithOptions returns a cache with options.
func NewServiceWithOptions(ctx context.Context, o *Options) (common.Service, error) {

	logger := log.FromContext(ctx)

	// create a Dapr service (e.g. ":50001", "0.0.0.0:50001", "10.1.1.1:50001" )
	s, err := daprd.NewService(o.Address)
	if err != nil {
		return nil, err
	}

	logger.Infof("Created dapr http service on address %v", o.Address)

	return s, nil
}

// NewService returns a cache.
func NewService(ctx context.Context) (common.Service, error) {

	o, err := NewOptions()
	if err != nil {
		return nil, err
	}

	return NewServiceWithOptions(ctx, o)
}
