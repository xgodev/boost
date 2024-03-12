package cloudevents

import (
	"sync"

	cloudevents "github.com/xgodev/boost/factory/contrib/cloudevents/sdk-go/v2"
	contextfx "github.com/xgodev/boost/factory/contrib/go.uber.org/fx/v1/module/core/context"
	"go.uber.org/fx"
)

type srvParams struct {
	fx.In
	Middlewares []Middleware `group:"_faas_middleware_"`
}

var handlerWrapperOnce sync.Once

// HandlerWrapperModule returns fx module for initialization of event handler wrapped in middleware.
//
// The module is only loaded once.
func HandlerWrapperModule() fx.Option {
	options := fx.Options()

	handlerWrapperOnce.Do(func() {

		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				DefaultHandlerWrapperOptions,

				func(handler cloudevents.Handler, options *HandlerWrapperOptions, m srvParams) *HandlerWrapper {
					return NewHandlerWrapper(handler, options, m.Middlewares...)
				},
			),
		)
	})

	return options
}

var helperOnce sync.Once

// HelperModule returns fx module for initialization of helper to start HTTP client for handlers.
//
// The module is only loaded once.
func HelperModule(extraOptions fx.Option) fx.Option {
	options := fx.Options()

	helperOnce.Do(func() {

		options = fx.Options(
			contextfx.Module(),
			extraOptions,
			HandlerWrapperModule(),
			fx.Provide(
				NewHelper,
			),
			fx.Invoke(
				func(helper *Helper) {
					helper.Start()
				},
			),
		)
	})

	return options
}
