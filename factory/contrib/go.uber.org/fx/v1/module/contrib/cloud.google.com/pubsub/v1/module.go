package bigquery

import (
	"sync"

	contextfx "github.com/xgodev/boost/factory/go.uber.org/fx.v1/module/context"
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
