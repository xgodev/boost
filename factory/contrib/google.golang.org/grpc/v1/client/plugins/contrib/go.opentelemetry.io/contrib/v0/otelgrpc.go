package contrib

import (
	"context"

	otelboost "github.com/xgodev/boost/factory/contrib/go.opentelemetry.io/otel/v1"
	"github.com/xgodev/boost/wrapper/log"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/stats/opentelemetry"
)

// allExperimentalMetrics lists every experimental metric registered in the
// gRPC-Go experimental/stats registry. They default to false upstream; this set
// opts into all of them.
var allExperimentalMetrics = stats.NewMetricSet(
	// grpc.subchannel.*
	"grpc.subchannel.disconnections",
	"grpc.subchannel.connection_attempts_succeeded",
	"grpc.subchannel.connection_attempts_failed",
	"grpc.subchannel.open_connections",
	// grpc.lb.pick_first.*
	"grpc.lb.pick_first.disconnections",
	"grpc.lb.pick_first.connection_attempts_succeeded",
	"grpc.lb.pick_first.connection_attempts_failed",
	// grpc.lb.rls.*
	"grpc.lb.rls.cache_entries",
	"grpc.lb.rls.cache_size",
	"grpc.lb.rls.default_target_picks",
	"grpc.lb.rls.target_picks",
	"grpc.lb.rls.failed_picks",
	// grpc.lb.wrr.*
	"grpc.lb.wrr.rr_fallback",
	"grpc.lb.wrr.endpoint_weight_not_yet_usable",
	"grpc.lb.wrr.endpoint_weight_stale",
	"grpc.lb.wrr.endpoint_weights",
	// grpc.lb.outlier_detection.*
	"grpc.lb.outlier_detection.ejections_enforced",
	"grpc.lb.outlier_detection.ejections_unenforced",
	// grpc.xds_client.*
	"grpc.xds_client.resource_updates_valid",
	"grpc.xds_client.resource_updates_invalid",
	"grpc.xds_client.server_failure",
	"grpc.xds_client.connected",
	"grpc.xds_client.resources",
)

// Register returns the gRPC DialOptions to enable OpenTelemetry stats handler.
// It reads the Enabled flag from Options; if disabled or on error, returns nils.
func Register(ctx context.Context) ([]grpc.DialOption, []grpc.CallOption) {
	o, err := NewOptions()
	if err != nil {
		return nil, nil
	}
	p := NewOpenTelemetryWithOptions(o)
	return p.Register(ctx)
}

// OpenTelemetry is the plugin for adding OpenTelemetry instrumentation to gRPC clients.
type OpenTelemetry struct {
	options *Options
}

// NewOpenTelemetryWithOptions creates a new plugin instance with the given options.
func NewOpenTelemetryWithOptions(options *Options) *OpenTelemetry {
	return &OpenTelemetry{options: options}
}

// NewOpenTelemetry creates a new plugin using default config.
func NewOpenTelemetry() *OpenTelemetry {
	o, err := NewOptions()
	if err != nil {
		log.Fatalf(err.Error())
	}
	return NewOpenTelemetryWithOptions(o)
}

// Register returns the DialOptions for OpenTelemetry. Uses StatsHandler instead of interceptors.
func (p *OpenTelemetry) Register(ctx context.Context) ([]grpc.DialOption, []grpc.CallOption) {
	if !p.options.Enabled {
		return nil, nil
	}

	logger := log.FromContext(ctx)
	logger.Debug("OpenTelemetry gRPC stats handler enabled")

	return []grpc.DialOption{
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		opentelemetry.DialOption(opentelemetry.Options{
			MetricsOptions: opentelemetry.MetricsOptions{
				MeterProvider: otelboost.MeterProvider,
				Metrics:       opentelemetry.DefaultMetrics().Join(allExperimentalMetrics),
			},
		}),
	}, nil
}
