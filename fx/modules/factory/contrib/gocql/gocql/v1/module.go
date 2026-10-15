package gocql

import (
	"context"
	"sync"

	"github.com/gocql/gocql"
	g "github.com/xgodev/boost/factory/contrib/gocql/gocql/v1"
	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"go.uber.org/fx"
)

var once sync.Once

type params struct {
	fx.In
	Plugins []g.Plugin `optional:"true"`
}

// Module fx module for gocql session. The session is closed on app stop.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				func(ctx context.Context, lc fx.Lifecycle, p params) (*gocql.Session, error) {
					s, err := g.NewSession(ctx, p.Plugins...)
					if err != nil {
						return nil, err
					}
					lc.Append(fx.StopHook(s.Close))
					return s, nil
				},
			),
		)
	})

	return options
}
