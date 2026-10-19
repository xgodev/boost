package otel

import (
	"context"

	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5"
	otelboost "github.com/xgodev/boost/factory/contrib/go.opentelemetry.io/otel/v1"
	"github.com/xgodev/boost/wrapper/log"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Register registers a new opentelemetry plugin for chi server.
func Register(ctx context.Context, server *chi.Server) error {
	h, err := NewOtel()
	if err != nil {
		return err
	}
	return h.Register(ctx, server)
}

// Otel represents opentelemetry plugin for chi server.
type Otel struct {
	options *Options
}

// NewOtelWithOptions returns a new opentelemetry plugin with options.
func NewOtelWithOptions(options *Options) *Otel {
	return &Otel{options: options}
}

// NewOtelWithConfigPath returns a new opentelemetry plugin with options from config path.
func NewOtelWithConfigPath(path string) (*Otel, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewOtelWithOptions(o), nil
}

// NewOtel returns a new opentelemetry plugin with default options.
func NewOtel() (*Otel, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewOtelWithOptions(o), nil
}

// Register registers this opentelemetry plugin for chi server. Spans are
// named "METHOD /route/{param}": chi sets the matched pattern on the request
// and otelhttp renames the span from it once the route ran.
func (i *Otel) Register(ctx context.Context, server *chi.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling opentelemetry middleware in chi")

	opts := []otelhttp.Option{
		otelhttp.WithTracerProvider(otelboost.TracerProvider),
		otelhttp.WithMeterProvider(otelboost.MeterProvider),
	}
	if i.options.ServerName != "" {
		opts = append(opts, otelhttp.WithServerName(i.options.ServerName))
	}
	server.Use(otelhttp.NewMiddleware(i.options.ServerName, opts...))

	logger.Debug("opentelemetry middleware successfully enabled in chi")

	return nil
}
