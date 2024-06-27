package publisher

import (
	"github.com/xgodev/boost/bootstrap/cloudevents"
	"sync"

	"github.com/xgodev/boost/bootstrap/wrapper/provider"
	"go.uber.org/fx"
)

var once sync.Once

// Module returns fx module for initialization of event publisher middleware.
// Which depends on event wrapper provider module.
//
// The module is only loaded once.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			provider.Module(),
			fx.Provide(
				NewOptions,
				fx.Annotated{
					Group: "_faas_middleware_",
					Target: func(options *Options, events *provider.EventWrapperProvider) cloudevents.Middleware {
						return NewEventPublisher(options, events)
					},
				},
			),
		)
	})

	return options
}
