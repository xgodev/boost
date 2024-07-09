package nats

import (
	fxnats "github.com/xgodev/boost/fx/modules/contrib/nats-io/nats.go/v1"
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
			fxnats.Module(),
			fx.Provide(
				nats.New,
			),
		)
	})

	return options
}
