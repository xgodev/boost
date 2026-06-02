package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub/v2"
	apiv1 "github.com/xgodev/boost/factory/contrib/cloud.google.com/api/v0"
	grpcv1 "github.com/xgodev/boost/factory/contrib/cloud.google.com/grpc/v1"
	otelboost "github.com/xgodev/boost/factory/contrib/go.opentelemetry.io/otel/v1"
	clientgrpc "github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/client"
	"github.com/xgodev/boost/wrapper/log"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats/opentelemetry"

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

	apiOpts := apiv1.ApplyAPIOptions(ctx, &o.APIOptions)
	grpcDialOpts := grpcv1.ApplyDialOptions(ctx, &o.GRPCOptions, plugins...)

	clientConfig := &pubsub.ClientConfig{}
	if o.Otel.Enabled {
		otelboost.StartMeterProvider(ctx)
		otelboost.StartTracerProvider(ctx)

		opts := opentelemetry.Options{
			MetricsOptions: opentelemetry.MetricsOptions{

				MeterProvider: otelboost.MeterProvider,
				Metrics:       opentelemetry.DefaultMetrics(),
			},
		}

		grpcDialOpts = append(grpcDialOpts, grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
		grpcDialOpts = append(grpcDialOpts, opentelemetry.DialOption(opts))
		clientConfig.EnableOpenTelemetryTracing = true
	}

	clientOpts := make([]option.ClientOption, 0, len(apiOpts)+len(grpcDialOpts))
	clientOpts = append(clientOpts, apiOpts...)
	for _, dop := range grpcDialOpts {
		clientOpts = append(clientOpts, option.WithGRPCDialOption(dop))
	}

	logger.Debugf("creating Pub/Sub client for project %s", o.APIOptions.ProjectID)

	return pubsub.NewClientWithConfig(ctx, o.APIOptions.ProjectID, clientConfig, clientOpts...)
}
