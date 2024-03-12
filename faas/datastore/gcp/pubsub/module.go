package pubsub

import (
	"sync"

	pubsubfx "github.com/xgodev/boost/factory/contrib/go.uber.org/fx/v1/module/contrib/cloud.google.com/pubsub/v1"
	"github.com/xgodev/boost/factory/contrib/go.uber.org/fx/v1/module/core/context"
	"go.uber.org/fx"
)

var once sync.Once

// Module loads the sns module providing an initialized client.
//
// The module is only loaded once.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			context.Module(),
			pubsubfx.Module(),
			fx.Provide(
				NewEvent,
			),
		)
	})

	return options
}
