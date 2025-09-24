package pubsub

import (
	"sync"

	"github.com/xgodev/boost/factory/contrib/cloud.google.com/pubsub/v2"
	config "github.com/xgodev/boost/fx/modules/factory/contrib/cloud.google.com/pubsub"

	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"go.uber.org/fx"
)

var once sync.Once

// Module fx module for pubsub client.
func Module() fx.Option {
	options := fx.Options()
	if config.ClientVersion() == "v2" {

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

	return options

}
