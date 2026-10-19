package kafka

import (
	"sync"

	"github.com/xgodev/boost/bootstrap/function/adapter/contrib/segmentio/kafka-go/v0"
	"github.com/xgodev/boost/fx/modules/bootstrap/function"
	"go.uber.org/fx"
)

var once sync.Once

// Module adds the adapter that runs the function on the messages of a kafka
// consumer group, on segmentio/kafka-go (the kafkago command), to the
// function adapters.
func Module[T any]() fx.Option {
	options := fx.Options()
	if !IsEnabled() {
		return options
	}

	once.Do(func() {
		options = fx.Options(
			fx.Provide(
				fx.Annotated{
					Group:  function.BSFunctionAdaptersGroupKey,
					Target: kafka.New[T],
				},
			),
		)
	})

	return options
}
