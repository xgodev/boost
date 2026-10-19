// Package elastic provides the fx modules of boost's go-elasticsearch v9
// factory.
package elastic

import (
	"context"
	"sync"

	"github.com/elastic/go-elasticsearch/v9"
	e "github.com/xgodev/boost/factory/contrib/elastic/go-elasticsearch/v9"
	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"go.uber.org/fx"
)

var once sync.Once

type clientParams struct {
	fx.In
	Plugins []e.Plugin `optional:"true"`
}

// Module provides the *elasticsearch.Client built from config, with the
// []e.Plugin provided in the graph. The client pings on creation, so an
// unreachable cluster fails the boot.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				func(ctx context.Context, p clientParams) (*elasticsearch.Client, error) {
					return e.NewClient(ctx, p.Plugins...)
				},
			),
		)
	})

	return options
}

var bonce sync.Once

// BulkModule provides the esutil.BulkIndexer built from config on the
// Module client.
func BulkModule() fx.Option {
	options := fx.Options()

	bonce.Do(func() {
		options = fx.Options(
			Module(),
			fx.Provide(
				e.NewBulkIndexer,
			),
		)
	})

	return options
}
