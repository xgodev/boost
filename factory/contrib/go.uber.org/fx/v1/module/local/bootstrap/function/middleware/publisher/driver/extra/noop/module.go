package noop

import (
	"github.com/xgodev/boost/bootstrap/function/middleware/publisher/driver/extra/noop"
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
				noop.New,
			),
		)
	})

	return options
}