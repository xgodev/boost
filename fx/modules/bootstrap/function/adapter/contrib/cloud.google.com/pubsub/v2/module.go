package pubsub

import (
	adaptercfg "github.com/xgodev/boost/fx/modules/bootstrap/function/adapter/contrib/cloud.google.com/pubsub"
	pubsubcfg "github.com/xgodev/boost/fx/modules/factory/contrib/cloud.google.com/pubsub"
	"sync"

	pubsub "github.com/xgodev/boost/bootstrap/function/adapter/contrib/cloud.google.com/pubsub/v2"
	"github.com/xgodev/boost/fx/modules/bootstrap/function"
	fxpubsub "github.com/xgodev/boost/fx/modules/factory/contrib/cloud.google.com/pubsub/v2"
	"go.uber.org/fx"
)

var once sync.Once

func Module[T any]() fx.Option {
	if _, err := pubsubcfg.Version(adaptercfg.VersionKey); err != nil {
		return fx.Error(err)
	}

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
