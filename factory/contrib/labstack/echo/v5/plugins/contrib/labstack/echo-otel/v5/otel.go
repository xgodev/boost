package otel

import (
	"context"
	echootel "github.com/labstack/echo-otel/v5"
	otelboost "github.com/xgodev/boost/factory/contrib/go.opentelemetry.io/otel/v1"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v5"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new opentelemetry plugin for echo server.
func Register(ctx context.Context, server *echo.Server) error {
	h, err := NewOtel()
	if err != nil {
		return err
	}
	return h.Register(ctx, server)
}

// Otel represents opentelemetry plugin for echo server.
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

// Register registers this opentelemetry plugin for echo server.
func (i *Otel) Register(ctx context.Context, server *echo.Server) error {
	if !i.options.Enabled {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling opentelemetry middleware in echo")

	mw, err := echootel.Config{
		ServerName:     i.options.ServerName,
		TracerProvider: otelboost.TracerProvider,
		MeterProvider:  otelboost.MeterProvider,
	}.ToMiddleware()
	if err != nil {
		return err
	}
	server.Use(mw)

	logger.Debug("opentelemetry middleware successfully enabled in echo")

	return nil
}
