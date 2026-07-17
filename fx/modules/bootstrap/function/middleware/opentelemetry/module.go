package opentelemetry

import (
	"sync"

	xotel "github.com/xgodev/boost/bootstrap/function/middleware/go.opentelemetry.io/otel/v1"
	"go.uber.org/fx"
)

var once sync.Once

func Module[T any]() fx.Option {
	options := fx.Options()
	if !IsEnabled() {
		return options
	}

	once.Do(func() {
		options = fx.Options(
			fx.Provide(
				xotel.NewOpenTelemetry[T],
			),
		)
	})

	return options
}
