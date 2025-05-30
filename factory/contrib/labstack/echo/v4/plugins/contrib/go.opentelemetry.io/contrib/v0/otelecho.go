package contrib

import (
	"context"
	"github.com/xgodev/boost/factory/contrib/go.opentelemetry.io/otel/v1"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4"

	"github.com/xgodev/boost/wrapper/log"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

// Register registers a new opentelemetry plugin for echo server.
func Register(ctx context.Context, server *echo.Server) error {
	o, err := NewOptions()
	if err != nil {
		return nil
	}
	h := NewOtelEchoWithOptions(o)
	h.Register(ctx, server)
	return nil
}

// OtelEcho represents opentelemetry plugin for echo server.
type OtelEcho struct {
	options *Options
}

// NewOtelEchoWithOptions returns a new opentelemetry plugin with options.
func NewOtelEchoWithOptions(options *Options) *OtelEcho {
	return &OtelEcho{options: options}
}

// NewOtelEchoWithConfigPath returns a new opentelemetry plugin with options from config path.
func NewOtelEchoWithConfigPath(path string, tracingOptions ...otelecho.Option) (*OtelEcho, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	o.TracingOptions = tracingOptions
	return NewOtelEchoWithOptions(o), nil
}

// NewOtelEcho returns a new opentelemetry plugin with default options.
func NewOtelEcho(tracingOptions ...otelecho.Option) *OtelEcho {
	o, err := NewOptions()
	if err != nil {
		log.Fatalf(err.Error())
	}

	o.TracingOptions = tracingOptions
	return NewOtelEchoWithOptions(o)
}

// Register registers this opentelemetry plugin for echo server.
func (i *OtelEcho) Register(ctx context.Context, server *echo.Server) {
	if !i.options.Enabled {
		return
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling opentelemetry middleware in echo")

	otel.StartTracerProvider(ctx)
	otel.StartMeterProvider(ctx)

	i.options.TracingOptions = append(i.options.TracingOptions, otelecho.WithTracerProvider(otel.TracerProvider))

	server.Use(otelecho.Middleware("", i.options.TracingOptions...))

	logger.Debug("opentelemetry integration successfully enabled in echo")
}
