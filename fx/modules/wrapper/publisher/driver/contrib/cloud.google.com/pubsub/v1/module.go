package pubsub

import (
	pubsubcfg "github.com/xgodev/boost/fx/modules/factory/contrib/cloud.google.com/pubsub"
	fxpubsub "github.com/xgodev/boost/fx/modules/factory/contrib/cloud.google.com/pubsub/v1"
	drivercfg "github.com/xgodev/boost/fx/modules/wrapper/publisher/driver/contrib/cloud.google.com/pubsub"
	"github.com/xgodev/boost/wrapper/publisher/driver/contrib/cloud.google.com/pubsub/v1"
	"go.uber.org/fx"
	"sync"
)

var once sync.Once

func Module() fx.Option {
	if _, err := pubsubcfg.Version(drivercfg.VersionKey); err != nil {
		return fx.Error(err)
	}

	options := fx.Options()
	if !IsEnabled() {
		return options
	}

	once.Do(func() {
		options = fx.Options(
			fxpubsub.Module(),
			fx.Provide(
				pubsub.New,
			),
		)
	})

	return options
}
