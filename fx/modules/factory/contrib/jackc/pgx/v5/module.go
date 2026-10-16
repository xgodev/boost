package pgx

import (
	"context"
	"database/sql"
	"sync"

	d "github.com/xgodev/boost/factory/contrib/jackc/pgx/v5"
	sqlf "github.com/xgodev/boost/factory/core/database/sql"
	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"go.uber.org/fx"
)

var once sync.Once

type params struct {
	fx.In
	Plugins []sqlf.Plugin `optional:"true"`
}

// Module fx module for pgx *sql.DB. The pool is closed on app stop.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				func(ctx context.Context, lc fx.Lifecycle, p params) (*sql.DB, error) {
					db, err := d.NewDB(ctx, p.Plugins...)
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
