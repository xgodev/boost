package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub/v2"
	clientgrpc "github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/client"
	"github.com/xgodev/boost/wrapper/log"
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
	logger.Debugf("creating Pub/Sub client for project %s", o.APIOptions.ProjectID)

	return pubsub.NewClient(ctx, o.APIOptions.ProjectID)
}
