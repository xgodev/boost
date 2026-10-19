package redis

import (
	"context"
	"github.com/xgodev/boost/factory/contrib/redis/go-redis/v9"
	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"sync"

	r "github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

var once sync.Once

type clusterParams struct {
	fx.In
	Plugins []redis.ClusterPlugin `optional:"true"`
}

// ClusterModule fx module for redis cluster client. The client closes on app stop.
func ClusterModule() fx.Option {
	options := fx.Options()

	once.Do(func() {

		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				func(ctx context.Context, lc fx.Lifecycle, p clusterParams) (*r.ClusterClient, error) {
					client, err := redis.NewClusterClient(ctx, p.Plugins...)
					if err != nil {
						return nil, err
					}
					lc.Append(fx.StopHook(client.Close))
					return client, nil
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

// ClientModule fx module for redis client. The client closes on app stop.
func ClientModule() fx.Option {
	options := fx.Options()

	once.Do(func() {

		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				func(ctx context.Context, lc fx.Lifecycle, p clientParams) (*r.Client, error) {
					client, err := redis.NewClient(ctx, p.Plugins...)
					if err != nil {
						return nil, err
					}
					lc.Append(fx.StopHook(client.Close))
					return client, nil
				},
			),
		)
	})

	return options
}
