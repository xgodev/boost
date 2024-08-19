package pubsub

import (
	"github.com/xgodev/boost/factory/contrib/cloud.google.com/pubsub/v1"
	"sync"

	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"go.uber.org/fx"
)

var once sync.Once

// Module fx module for bigQuery client.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {

		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				pubsub.NewClient,
			),
		)

	})

	return options
}
