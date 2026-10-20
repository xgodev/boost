package otel

import (
	"context"
	"sync"

	xotel "github.com/xgodev/boost/factory/contrib/go.opentelemetry.io/otel/v1"
	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"github.com/xgodev/boost/wrapper/log"
	"go.uber.org/fx"
)

var once sync.Once

// Module starts the OpenTelemetry tracer and meter providers. Include this
// module in any application that uses OTel-instrumented clients (PubSub, gRPC,
// MongoDB, etc.) so that providers are initialized before the clients are created.
// On stop it shuts both providers down within the stop context, exporting the
// spans and metrics they still hold; an export that fails is logged at ERROR
// and does not fail the stop.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			contextfx.Module(),
			fx.Invoke(func(ctx context.Context, lc fx.Lifecycle) {
				xotel.StartTracerProvider(ctx)
				xotel.StartMeterProvider(ctx)
				lc.Append(fx.StopHook(func(stopCtx context.Context) {
					shutdown(stopCtx, xotel.Shutdown)
				}))
			}),
		)
	})

	return options
}

// shutdown runs stop and logs its error instead of returning it: telemetry
// the collector did not take must not turn a clean exit (a job that
// succeeded) into exit code 1.
func shutdown(ctx context.Context, stop func(context.Context) error) {
	if err := stop(ctx); err != nil {
		log.FromContext(ctx).Errorf("opentelemetry: %v", err)
	}
}
