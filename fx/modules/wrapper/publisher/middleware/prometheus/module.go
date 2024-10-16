package noop

import (
	"github.com/xgodev/boost/wrapper/publisher/middleware/prometheus"
	"go.uber.org/fx"
	"sync"
)

var once sync.Once

func Module() fx.Option {
	options := fx.Options()
	if !IsEnabled() {
		return options
	}

	once.Do(func() {
		options = fx.Options(
			fx.Provide(
				prometheus.NewPrometheus,
			),
		)
	})

	return options
}
