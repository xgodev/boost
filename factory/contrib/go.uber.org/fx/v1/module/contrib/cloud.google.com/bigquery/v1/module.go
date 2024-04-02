package bigquery

import (
	"github.com/xgodev/boost/factory/contrib/cloud.google.com/bigquery/v1"
	"sync"

	contextfx "github.com/xgodev/boost/factory/contrib/go.uber.org/fx/v1/module/core/context"
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
				bigquery.NewClient,
			),
		)

	})

	return options
}
