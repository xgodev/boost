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
// received over NATS (the cloudevents_nats command) to the function adapters,
// with the *nats.Conn of the nats factory fx module.
func Module[T any]() fx.Option {
	return module[T](fxnats.Module)
}

// ModuleWithoutClient adds the CloudEvents over NATS adapter to the function
// adapters and takes its *nats.Conn from the app: the one the app's publisher
// and @Inject consumers share. builder wires it with the default client
// inject generates.
func ModuleWithoutClient[T any]() fx.Option {
	return module[T](nil)
}

// module returns the adapter, with the module of its client unless client is
// nil, once. client is called only then, so a disabled adapter leaves the
// factory module to the app.
func module[T any](client func() fx.Option) fx.Option {
	options := fx.Options()
	if !IsEnabled() {
		return options
	}

	once.Do(func() {
		options = fx.Provide(
			fx.Annotated{
				Group:  function.BSFunctionAdaptersGroupKey,
				Target: nats.New[T],
			},
		)
		if client != nil {
			options = fx.Options(client(), options)
		}
	})

	return options
}
