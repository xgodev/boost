package ignore_errors

import (
	p "github.com/xgodev/boost/bootstrap/function/middleware/prometheus"
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
			fx.Provide(
				p.NewPrometheus[T],
			),
		)
	})

	return options
}
