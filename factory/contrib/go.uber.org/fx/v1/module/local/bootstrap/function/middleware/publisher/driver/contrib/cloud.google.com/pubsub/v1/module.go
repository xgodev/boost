package pubsub

import (
	"github.com/xgodev/boost/bootstrap/function/middleware/publisher/driver/contrib/cloud.google.com/pubsub/v1"
	fxpubsub "github.com/xgodev/boost/factory/contrib/go.uber.org/fx/v1/module/contrib/cloud.google.com/pubsub/v1"
	"go.uber.org/fx"
	"sync"
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
			),
		)
	})

	return options
}
