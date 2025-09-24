package pubsub

import (
	"context"
	"sync"

	fxpubsub "github.com/xgodev/boost/fx/modules/factory/contrib/cloud.google.com/pubsub/v2"
	"github.com/xgodev/boost/wrapper/publisher"
	"github.com/xgodev/boost/wrapper/publisher/driver/contrib/cloud.google.com/pubsub/v2"
	"go.uber.org/fx"
)

var once sync.Once

func Module() fx.Option {
	options := fx.Options()
	if !IsEnabled() {
		return options
	}

	once.Do(func() {
		options = fx.Options(
			fxpubsub.Module(),
			fx.Provide(
				pubsub.New,
				shut,
			),
		)
	})

	return options
}

func shut(lc fx.Lifecycle, pub publisher.Driver) fx.Option {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			pub.Close()
			return nil
		},
	})

	return fx.Options()
}
