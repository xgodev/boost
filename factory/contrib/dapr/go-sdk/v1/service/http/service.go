package http

import (
	"context"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
	"github.com/xgodev/boost/wrapper/log"
)

// NewWithConfigPath returns a cache with options from config path .
func NewWithConfigPath(ctx context.Context, path string) (common.Service, error) {
	options, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewWithOptions(ctx, options), nil
}

// NewWithOptions returns a cache with options.
func NewWithOptions(ctx context.Context, o *Options) common.Service {

	logger := log.FromContext(ctx)

	// create a Dapr service (e.g. ":8080", "0.0.0.0:8080", "10.1.1.1:8080" )
	s := daprd.NewService(o.Address)

	logger.Infof("Created dapr http service on address %v", o.Address)

	return s
}

// NewService returns a cache.
func NewService(ctx context.Context) (common.Service, error) {

	o, err := NewOptions()
	if err != nil {
		return nil, err
	}

	return NewWithOptions(ctx, o), nil
}
