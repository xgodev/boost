package nats

import (
	"sync"

	"github.com/xgodev/boost/bootstrap/function/adapter/contrib/cloudevents/sdk-go/v2/contrib/nats-io/nats.go/v1"
	"github.com/xgodev/boost/fx/modules/bootstrap/function"
	fxnats "github.com/xgodev/boost/fx/modules/factory/contrib/nats-io/nats.go/v1"
	"go.uber.org/fx"
)

var once sync.Once

// Module adds the adapter that feeds the function with the CloudEvents
// received over NATS (the cloudevents_nats command) to the function adapters.
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
