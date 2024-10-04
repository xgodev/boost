package logger

import (
	"github.com/xgodev/boost/bootstrap/function/middleware/logger"
	"go.uber.org/fx"
	"sync"
)

var once sync.Once

func Module[T any]() fx.Option {
	options := fx.Options()
	if !IsEnabled() {
		return options
	}

	once.Do(func() {
		options = fx.Options(
			fx.Provide(logger.NewOptions),
			fx.Provide(
				logger.NewLogger[T],
			),
		)
	})

	return options
}
