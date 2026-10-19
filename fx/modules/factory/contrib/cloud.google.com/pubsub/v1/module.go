package pubsub

import (
	"sync"

	"github.com/xgodev/boost/factory/contrib/cloud.google.com/pubsub/v1"
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

	if !IsEnabled() {
		return fx.Options()
	}

	return ClientModule()
}

// ClientModule provides the v1 *pubsub.Client whatever
// boost.factory.gcp.pubsub.version says. Modules that read with the v1 client,
// like the v1 function adapter, include it. It shares Module's once, so the
// client is provided a single time when both are included.
func ClientModule() fx.Option {
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
