package bigquery

import (
	"context"

	"cloud.google.com/go/bigquery"
	apiv1 "github.com/xgodev/boost/factory/contrib/cloud.google.com/api/v0"
	grpcv1 "github.com/xgodev/boost/factory/contrib/cloud.google.com/grpc/v1"
	clientgrpc "github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/client"
	"github.com/xgodev/boost/wrapper/log"
	"google.golang.org/api/option"
)

// Client wraps a BigQuery client and its Options.
type Client struct {
	Inner *bigquery.Client
	Opts  *Options
}

// NewClient creates a wrapped BigQuery client using default configuration.
func NewClient(ctx context.Context, plugins ...clientgrpc.Plugin) (*bigquery.Client, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewClientWithOptions(ctx, o, plugins...)
}

// NewClientWithConfigPath creates a wrapped BigQuery client using configuration from the specified path.
func NewClientWithConfigPath(ctx context.Context, path string, plugins ...clientgrpc.Plugin) (*bigquery.Client, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewClientWithOptions(ctx, o, plugins...)
}

// NewClientWithOptions constructs a wrapped BigQuery client from Options.
func NewClientWithOptions(ctx context.Context, o *Options, plugins ...clientgrpc.Plugin) (*bigquery.Client, error) {
	logger := log.FromContext(ctx)

	// API-level options
	apiOpts := apiv1.ApplyAPIOptions(ctx, &o.APIOptions)

	// gRPC-level DialOptions
	grpcDialOpts := grpcv1.ApplyDialOptions(ctx, &o.GRPCOptions, plugins...)

	// collect ClientOption
	var clientOpts []option.ClientOption
	clientOpts = append(clientOpts, apiOpts...)
	for _, dop := range grpcDialOpts {
		clientOpts = append(clientOpts, option.WithGRPCDialOption(dop))
	}

	// Quota project override
	if o.UserProject != "" {
		clientOpts = append(clientOpts, option.WithQuotaProject(o.UserProject))
	}

	logger.Debugf("creating BigQuery client for project %s", o.APIOptions.ProjectID)
	bqClient, err := bigquery.NewClient(ctx, o.APIOptions.ProjectID, clientOpts...)
	if err != nil {
		return nil, err
	}
	return bqClient, nil
}
