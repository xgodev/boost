package pubsub

import (
	"sync"

	pubsub "github.com/xgodev/boost/factory/contrib/cloud.google.com/pubsub/v2"
	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	pubsubcfg "github.com/xgodev/boost/fx/modules/factory/contrib/cloud.google.com/pubsub"
	"go.uber.org/fx"
)

var once sync.Once

// Module fx module for pubsub client.
func Module() fx.Option {
	if _, err := pubsubcfg.Version(pubsubcfg.VersionKey); err != nil {
		return fx.Error(err)
	}

	options := fx.Options()
	if !IsEnabled() {
		return options
	}

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
