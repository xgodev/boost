package pubsub

import (
	"github.com/xgodev/boost/bootstrap/function/adapter/contrib/cloud.google.com/pubsub/v1"
	"github.com/xgodev/boost/fx/modules/bootstrap/function"
	fxpubsub "github.com/xgodev/boost/fx/modules/factory/contrib/cloud.google.com/pubsub/v1"
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
			fxpubsub.Module(),
			fx.Provide(
				fx.Annotated{
					Group:  function.BSFunctionAdaptersGroupKey,
					Target: pubsub.New[T],
				},
			),
		)
	})

	return options
}
