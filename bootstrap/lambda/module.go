package lambda

import (
	"sync"

	"github.com/xgodev/boost/bootstrap/cloudevents"
	contextfx "github.com/xgodev/boost/factory/contrib/go.uber.org/fx/v1/module/core/context"
	"go.uber.org/fx"
)

var once sync.Once

// HelperModule returns fx module for initialization of helper to start HTTP client for handlers.
//
// The module is only loaded once.
func HelperModule(extraOptions fx.Option) fx.Option {
	options := fx.Options()

	once.Do(func() {

		options = fx.Options(
			contextfx.Module(),
			extraOptions,
			cloudevents.HandlerWrapperModule(),
			fx.Provide(
				NewDefaultHelper,
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
