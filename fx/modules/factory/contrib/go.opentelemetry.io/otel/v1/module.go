package otel

import (
	"context"
	"sync"

	xotel "github.com/xgodev/boost/factory/contrib/go.opentelemetry.io/otel/v1"
	"go.uber.org/fx"
)

var once sync.Once

// Module starts the OpenTelemetry tracer and meter providers. Include this
// module in any application that uses OTel-instrumented clients (PubSub, gRPC,
// MongoDB, etc.) so that providers are initialized before the clients are created.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			fx.Invoke(func(ctx context.Context) {
				xotel.StartTracerProvider(ctx)
				xotel.StartMeterProvider(ctx)
			}),
		)
	})

	return options
}
