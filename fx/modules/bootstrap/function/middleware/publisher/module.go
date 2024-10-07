package publisher

import (
	p "github.com/xgodev/boost/bootstrap/function/middleware/publisher"
	"github.com/xgodev/boost/fx/modules/wrapper/publisher"
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
			publisher.Module(),
			fx.Provide(
				p.NewPublisher[T],
			),
		)
	})

	return options
}
