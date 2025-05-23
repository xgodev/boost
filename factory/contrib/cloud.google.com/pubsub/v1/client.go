package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub"
	apiv1 "github.com/xgodev/boost/factory/contrib/cloud.google.com/api/v0"
	grpcv1 "github.com/xgodev/boost/factory/contrib/cloud.google.com/grpc/v1"
	clientgrpc "github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/client"
	"github.com/xgodev/boost/wrapper/log"
	"google.golang.org/api/option"
)

// NewClient creates a Pub/Sub client using default configuration.
func NewClient(ctx context.Context, plugins ...clientgrpc.Plugin) (*pubsub.Client, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewClientWithOptions(ctx, o, plugins...)
}

// NewClientWithConfigPath creates a Pub/Sub client from a specific config path.
func NewClientWithConfigPath(ctx context.Context, path string, plugins ...clientgrpc.Plugin) (*pubsub.Client, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewClientWithOptions(ctx, o, plugins...)
}

// NewClientWithOptions constructs a Pub/Sub client from Options.
func NewClientWithOptions(ctx context.Context, o *Options, plugins ...clientgrpc.Plugin) (*pubsub.Client, error) {
	logger := log.FromContext(ctx)

	// API-level options
	apiOpts := apiv1.ApplyAPIOptions(ctx, &o.APIOptions)

	// gRPC-level DialOptions
	grpcDialOpts := grpcv1.ApplyDialOptions(ctx, &o.GRPCOptions, plugins...)

	// collect ClientOption
	clientOpts := make([]option.ClientOption, 0, len(apiOpts)+len(grpcDialOpts))
	clientOpts = append(clientOpts, apiOpts...)
	for _, dop := range grpcDialOpts {
		clientOpts = append(clientOpts, option.WithGRPCDialOption(dop))
	}

	logger.Debugf("creating Pub/Sub client for project %s", o.APIOptions.ProjectID)
	return pubsub.NewClient(ctx, o.APIOptions.ProjectID, clientOpts...)
}
