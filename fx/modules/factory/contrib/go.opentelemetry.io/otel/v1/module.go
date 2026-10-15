package otel

import (
	"context"
	"sync"

	o "github.com/xgodev/boost/factory/contrib/go.opentelemetry.io/otel/v1"
	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"go.uber.org/fx"
)

var optOnce sync.Once

// OptionsModule fx module for opentelemetry options.
func OptionsModule() fx.Option {
	options := fx.Options()

	optOnce.Do(func() {
		options = fx.Options(
			fx.Provide(
				o.NewOptions,
			),
		)
	})

	return options
}

var tracerOnce sync.Once

// TracerModule fx module that starts the opentelemetry tracer provider.
func TracerModule() fx.Option {
	options := fx.Options()

	tracerOnce.Do(func() {
		options = fx.Options(
			contextfx.Module(),
			OptionsModule(),
			fx.Invoke(
				func(ctx context.Context, opts *o.Options) {
					o.StartTracerProviderWithOptions(ctx, opts)
				},
			),
		)
	})

	return options
}

var meterOnce sync.Once

// MeterModule fx module that starts the opentelemetry meter provider.
func MeterModule() fx.Option {
	options := fx.Options()

	meterOnce.Do(func() {
		options = fx.Options(
			contextfx.Module(),
			OptionsModule(),
			fx.Invoke(
				func(ctx context.Context, opts *o.Options) {
					o.StartMetricProviderWithOptions(ctx, opts)
				},
			),
		)
	})

	return options
}
