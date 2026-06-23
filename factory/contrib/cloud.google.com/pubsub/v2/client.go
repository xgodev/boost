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
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/stats/opentelemetry"

	"google.golang.org/api/option"
)

// allExperimentalMetrics lists every experimental metric registered in the
// gRPC-Go experimental/stats registry.
var allExperimentalMetrics = stats.NewMetricSet(
	"grpc.subchannel.disconnections",
	"grpc.subchannel.connection_attempts_succeeded",
	"grpc.subchannel.connection_attempts_failed",
	"grpc.subchannel.open_connections",
	"grpc.lb.pick_first.disconnections",
	"grpc.lb.pick_first.connection_attempts_succeeded",
	"grpc.lb.pick_first.connection_attempts_failed",
	"grpc.lb.rls.cache_entries",
	"grpc.lb.rls.cache_size",
	"grpc.lb.rls.default_target_picks",
	"grpc.lb.rls.target_picks",
	"grpc.lb.rls.failed_picks",
	"grpc.lb.wrr.rr_fallback",
	"grpc.lb.wrr.endpoint_weight_not_yet_usable",
	"grpc.lb.wrr.endpoint_weight_stale",
	"grpc.lb.wrr.endpoint_weights",
	"grpc.lb.outlier_detection.ejections_enforced",
	"grpc.lb.outlier_detection.ejections_unenforced",
	"grpc.xds_client.resource_updates_valid",
	"grpc.xds_client.resource_updates_invalid",
	"grpc.xds_client.server_failure",
	"grpc.xds_client.connected",
	"grpc.xds_client.resources",
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
		opts := opentelemetry.Options{
			MetricsOptions: opentelemetry.MetricsOptions{
				MeterProvider: otelboost.MeterProvider,
				Metrics:       opentelemetry.DefaultMetrics().Join(allExperimentalMetrics),
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
