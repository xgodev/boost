package nats

import (
	fxnats "github.com/xgodev/boost/fx/modules/factory/contrib/nats-io/nats.go/v1"
	"github.com/xgodev/boost/wrapper/publisher/driver/contrib/nats-io/nats.go/v1"
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
