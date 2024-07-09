package nats

import (
	"github.com/xgodev/boost/bootstrap/function/middleware/publisher/driver/contrib/nats-io/nats.go/v1"
	fxnats "github.com/xgodev/boost/factory/contrib/go.uber.org/fx/v1/module/contrib/nats-io/nats.go/v1"
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
