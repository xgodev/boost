package contrib

import (
	"context"

	"github.com/xgodev/boost/wrapper/log"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
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

	// Use StatsHandler as recommended (UnaryClientInterceptor and StreamClientInterceptor are deprecated)
	return []grpc.DialOption{
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	}, nil
}
