package confluent

import (
	"github.com/xgodev/boost/bootstrap/function/adapter/contrib/confluentinc/confluent-kafka-go/v2"
	"github.com/xgodev/boost/fx/modules/bootstrap/function"
	fxconfluent "github.com/xgodev/boost/fx/modules/factory/contrib/confluentinc/confluent-kafka-go/v2"
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
			fxconfluent.ConsumerModule(),
			fx.Provide(
				fx.Annotated{
					Group:  function.BSFunctionAdaptersGroupKey,
					Target: confluent.New[T],
				},
			),
		)
	})

	return options
}
