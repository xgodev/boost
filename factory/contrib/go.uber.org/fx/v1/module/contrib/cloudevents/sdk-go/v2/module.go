package cloudevents

import (
	"sync"

	cloudevents "github.com/xgodev/boost/factory/contrib/cloudevents/sdk-go/v2"
	contextfx "github.com/xgodev/boost/factory/contrib/go.uber.org/fx/v1/module/core/context"
	"go.uber.org/fx"
)

var once sync.Once

// Module fx module for cloudevents.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {

		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				cloudevents.NewHTTP,
			),
		)
	})

	return options
}
