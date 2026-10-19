package datadog

import (
	"context"

	fibertrace "github.com/DataDog/dd-trace-go/contrib/gofiber/fiber.v2/v2"
	datadog "github.com/xgodev/boost/factory/contrib/datadog/dd-trace-go/v1"
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
	"github.com/xgodev/boost/wrapper/log"
)

// Register registers a new datadog middleware plugin for fiber server.
func Register(ctx context.Context, server *fiber.Server) error {
	p, err := NewDatadog()
	if err != nil {
		return err
	}
	return p.Register(ctx, server)
}

// Datadog represents datadog middleware plugin for fiber server.
type Datadog struct {
	options *Options
}

// NewDatadogWithOptions returns a new datadog middleware plugin with options.
func NewDatadogWithOptions(options *Options) *Datadog {
	return &Datadog{options: options}
}

// NewDatadogWithConfigPath returns a new datadog middleware plugin with options from config path.
func NewDatadogWithConfigPath(path string, traceOptions ...fibertrace.Option) (*Datadog, error) {
	o, err := NewOptionsWithPath(path, traceOptions...)
	if err != nil {
		return nil, err
	}
	return NewDatadogWithOptions(o), nil
}

// NewDatadog returns a new datadog middleware plugin with default options.
func NewDatadog(traceOptions ...fibertrace.Option) (*Datadog, error) {
	o, err := NewOptions(traceOptions...)
	if err != nil {
		return nil, err
	}
	return NewDatadogWithOptions(o), nil
}

// Register registers this datadog plugin for fiber server. It does nothing
// while the datadog tracer is disabled.
func (i *Datadog) Register(ctx context.Context, server *fiber.Server) error {
	if !i.options.Enabled || !datadog.IsTracerEnabled() {
		return nil
	}

	logger := log.FromContext(ctx)

	logger.Trace("enabling datadog middleware in fiber")

	server.Use(fibertrace.Middleware(i.options.TraceOptions...))

	logger.Debug("datadog middleware successfully enabled in fiber")

	return nil
}
