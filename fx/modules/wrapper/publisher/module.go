package publisher

import (
	"github.com/xgodev/boost/wrapper/publisher"
	"go.uber.org/fx"
	"sync"
)

var once sync.Once

func Module() fx.Option {

	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			fx.Provide(
				publisher.New,
			),
		)
	})

	return options
}
