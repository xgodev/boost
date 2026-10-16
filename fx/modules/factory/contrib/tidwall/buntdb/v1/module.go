package buntdb

import (
	"context"
	"sync"

	"github.com/tidwall/buntdb"
	b "github.com/xgodev/boost/factory/contrib/tidwall/buntdb/v1"
	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"go.uber.org/fx"
)

var once sync.Once

// Module fx module for buntdb. The database is closed on app stop.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				func(ctx context.Context, lc fx.Lifecycle) (*buntdb.DB, error) {
					db, err := b.NewDB(ctx)
					if err != nil {
						return nil, err
					}
					lc.Append(fx.StopHook(db.Close))
					return db, nil
				},
			),
		)
	})

	return options
}
