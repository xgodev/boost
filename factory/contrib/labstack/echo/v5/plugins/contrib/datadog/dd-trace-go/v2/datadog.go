package datadog

import (
	"context"

	echotrace "github.com/DataDog/dd-trace-go/contrib/labstack/echo.v5/v2"
	datadog "github.com/xgodev/boost/factory/contrib/datadog/dd-trace-go/v1"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v5"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new datadog plugin for echo server.
func Register(ctx context.Context, server *echo.Server) error {
	h, err := NewDatadog()
	if err != nil {
		return err
	}
	return h.Register(ctx, server)
}

// Datadog represents datadog plugin for echo server.
type Datadog struct {
	options *Options
}

// NewDatadogWithOptions returns a new datadog plugin with options.
func NewDatadogWithOptions(options *Options) *Datadog {
	return &Datadog{options: options}
}

// NewDatadogWithConfigPath returns a new datadog plugin with options from config path.
func NewDatadogWithConfigPath(path string, traceOptions ...echotrace.Option) (*Datadog, error) {
	o, err := NewOptionsWithPath(path, traceOptions...)
	if err != nil {
		return nil, err
	}
	return NewDatadogWithOptions(o), nil
}

// NewDatadog returns a new datadog plugin with default options.
func NewDatadog(traceOptions ...echotrace.Option) (*Datadog, error) {
	o, err := NewOptions(traceOptions...)
	if err != nil {
		return nil, err
	}
	return NewDatadogWithOptions(o), nil
}

// Register registers this datadog plugin for echo server. It does nothing
// while the datadog tracer is disabled.
func (i *Datadog) Register(ctx context.Context, server *echo.Server) error {
	if !i.options.Enabled || !datadog.IsTracerEnabled() {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling datadog middleware in echo")

	server.Use(echotrace.Middleware(i.options.TraceOptions...))

	logger.Debug("datadog middleware successfully enabled in echo")

	return nil
}
