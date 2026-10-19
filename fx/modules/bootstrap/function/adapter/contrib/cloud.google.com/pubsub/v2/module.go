package pubsub

import (
	"sync"

	pubsub "github.com/xgodev/boost/bootstrap/function/adapter/contrib/cloud.google.com/pubsub/v2"
	"github.com/xgodev/boost/fx/modules/bootstrap/function"
	adaptercfg "github.com/xgodev/boost/fx/modules/bootstrap/function/adapter/contrib/cloud.google.com/pubsub"
	pubsubcfg "github.com/xgodev/boost/fx/modules/factory/contrib/cloud.google.com/pubsub"
	fxpubsub "github.com/xgodev/boost/fx/modules/factory/contrib/cloud.google.com/pubsub/v2"
	"go.uber.org/fx"
)

var once sync.Once

// Module adds the Pub/Sub v2 adapter to the function adapters, with the
// *pubsub.Client of the pubsub v2 factory fx module. It is empty unless the
// adapter is enabled and its version is v2.
func Module[T any]() fx.Option {
	return module[T](fxpubsub.Module)
}

// ModuleWithoutClient adds the Pub/Sub v2 adapter to the function adapters
// and takes its *pubsub.Client from the app: the one the app's publisher
// and @Inject consumers share. builder wires it with the default client
// inject generates. It is empty unless the adapter is enabled and its
// version is v2.
func ModuleWithoutClient[T any]() fx.Option {
	return module[T](nil)
}

// ModuleWithNamedClient adds the Pub/Sub v2 adapter to the function adapters
// and takes its *pubsub.Client from the app: the one named name, which
// builder wires for @Function(client=name) with the named client inject
// generates. It is empty unless the adapter is enabled and its version is v2.
func ModuleWithNamedClient[T any](name string) fx.Option {
	return moduleOf(func() fx.Option {
		return fx.Provide(fx.Annotate(pubsub.New[T],
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
				Target: pubsub.New[T],
			},
		)
	}, client)
}

// moduleOf returns the option adapter provides, with the module of its
// client unless client is nil, once.
func moduleOf(adapter, client func() fx.Option) fx.Option {
	if _, err := pubsubcfg.Version(adaptercfg.VersionKey); err != nil {
		return fx.Error(err)
	}

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
