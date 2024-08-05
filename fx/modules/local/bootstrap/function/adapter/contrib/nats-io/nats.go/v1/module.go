package nats

import (
	"github.com/xgodev/boost/bootstrap/function/adapter/contrib/nats-io/nats.go/v1"
	fxnats "github.com/xgodev/boost/fx/modules/contrib/nats-io/nats.go/v1"
	"github.com/xgodev/boost/fx/modules/local/bootstrap/function"
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
			fxnats.Module(),
			fx.Provide(
				fx.Annotated{
					Group:  function.BSFunctionAdaptersGroupKey,
					Target: nats.New[T],
				},
			),
		)
	})

	return options
}
