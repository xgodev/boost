package mongo

import (
	"context"
	"github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v1"
	fxcontext "github.com/xgodev/boost/factory/contrib/go.uber.org/fx/v1/module/core/context"
	m "go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"sync"
)

const (
	PluginsGroupKey = "factory.mongo.plugins"
)

type params struct {
	fx.In
	Plugins []mongo.Plugin `optional:"true"`
}

var once sync.Once

func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Module("mongo",
			fxcontext.Module(),
			fx.Provide(
				func(ctx context.Context, p params) (*mongo.Conn, error) {
					return mongo.NewConn(ctx, p.Plugins...) //nolint:wrapcheck  // This must be resolved in another task
				},
				func(conn *mongo.Conn) *m.Database {
					return conn.Database
				},
			),
		)
	})

	return options
}
