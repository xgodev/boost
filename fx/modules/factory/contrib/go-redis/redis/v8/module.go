package redis

import (
	"context"
	"github.com/xgodev/boost/factory/contrib/go-redis/redis/v8"
	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"sync"

	r "github.com/go-redis/redis/v8"
	"go.uber.org/fx"
)

var once sync.Once

type clusterParams struct {
	fx.In
	Plugins []redis.ClusterPlugin `optional:"true"`
}

// ClusterModule fx module for redis cluster client.
func ClusterModule() fx.Option {
	options := fx.Options()

	once.Do(func() {

		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				func(ctx context.Context, p clusterParams) (*r.ClusterClient, error) {
					return redis.NewClusterClient(ctx, p.Plugins...)
				},
			),
		)
	})

	return options
}

type clientParams struct {
	fx.In
	Plugins []redis.Plugin `optional:"true"`
}

// ClientModule fx module for redis client.
func ClientModule() fx.Option {
	options := fx.Options()

	once.Do(func() {

		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				func(ctx context.Context, p clientParams) (*r.Client, error) {
					return redis.NewClient(ctx, p.Plugins...)
				},
			),
		)
	})

	return options
}
