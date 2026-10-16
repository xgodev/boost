package memdb

import (
	"context"
	"sync"

	"github.com/hashicorp/go-memdb"
	m "github.com/xgodev/boost/factory/contrib/hashicorp/memdb/v1"
	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"go.uber.org/fx"
)

var once sync.Once

// Module fx module for go-memdb. The *memdb.DBSchema must be provided by the app.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				func(ctx context.Context, schema *memdb.DBSchema) (*memdb.MemDB, error) {
					return m.NewMemDB(ctx, schema)
				},
			),
		)
	})

	return options
}
