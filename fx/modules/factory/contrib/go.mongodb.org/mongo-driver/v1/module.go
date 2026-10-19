package mongo

import (
	"context"
	"sync"

	"github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v1"
	fxcontext "github.com/xgodev/boost/fx/modules/core/context"
	m "go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
)

const (
	PluginsGroupKey = "factory.mongo.plugins"
)

type params struct {
	fx.In
	Plugins []mongo.Plugin `optional:"true"`
}

var once sync.Once

// Module fx module for the mongo *Conn. The client disconnects on app stop.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Module("mongo",
			fxcontext.Module(),
			fx.Provide(
				func(ctx context.Context, lc fx.Lifecycle, p params) (*mongo.Conn, error) {
					conn, err := mongo.NewConn(ctx, p.Plugins...)
					if err != nil {
						return nil, err //nolint:wrapcheck  // This must be resolved in another task
					}
					lc.Append(fx.StopHook(conn.Close))
					return conn, nil
				},
				func(conn *mongo.Conn) *m.Database {
					return conn.Database
				},
			),
		)
	})

	return options
}
