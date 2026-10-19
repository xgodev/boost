package nats

import (
	"sync"

	"github.com/xgodev/boost/bootstrap/function/adapter/contrib/nats-io/nats.go/v1"
	"github.com/xgodev/boost/fx/modules/bootstrap/function"
	fxnats "github.com/xgodev/boost/fx/modules/factory/contrib/nats-io/nats.go/v1"
	"go.uber.org/fx"
)

var once sync.Once

// Module adds the nats adapter to the function adapters, with the *nats.Conn
// of the nats factory fx module.
func Module[T any]() fx.Option {
	return module[T](fxnats.Module)
}

// ModuleWithoutClient adds the nats adapter to the function adapters and
// takes its *nats.Conn from the app: the one the app's publisher and @Inject
// consumers share. builder wires it with the default client inject generates.
func ModuleWithoutClient[T any]() fx.Option {
	return module[T](nil)
}

// ModuleWithNamedClient adds the nats adapter to the function adapters and
// takes its *nats.Conn from the app: the one named name, which builder wires
// for @Function(client=name) with the named client inject generates.
func ModuleWithNamedClient[T any](name string) fx.Option {
	return moduleOf(func() fx.Option {
		return fx.Provide(fx.Annotate(nats.New[T],
			fx.ParamTags(`name:"`+name+`"`),
			fx.ResultTags(`group:"`+function.BSFunctionAdaptersGroupKey+`"`),
		))
	}, nil)
}

// module returns the adapter, with the module of its client unless client is
// nil, once. client is called only then, so a disabled adapter leaves the
// factory module to the app.
func module[T any](client func() fx.Option) fx.Option {
	return moduleOf(func() fx.Option {
		return fx.Provide(
			fx.Annotated{
				Group:  function.BSFunctionAdaptersGroupKey,
				Target: nats.New[T],
			},
		)
	}, client)
}

// moduleOf returns the option adapter provides, with the module of its
// client unless client is nil, once.
func moduleOf(adapter, client func() fx.Option) fx.Option {
	options := fx.Options()
	if !IsEnabled() {
		return options
	}

	once.Do(func() {
		options = adapter()
		if client != nil {
			options = fx.Options(client(), options)
		}
	})

	return options
}
