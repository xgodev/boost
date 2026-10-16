package bigcache

import (
	"context"
	"sync"

	"github.com/allegro/bigcache/v3"
	b "github.com/xgodev/boost/factory/contrib/allegro/bigcache/v3"
	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"go.uber.org/fx"
)

var once sync.Once

// Module fx module for bigcache. The cache is closed on app stop.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				func(ctx context.Context, lc fx.Lifecycle) (*bigcache.BigCache, error) {
					cache, err := b.NewCache(ctx)
					if err != nil {
						return nil, err
					}
					lc.Append(fx.StopHook(cache.Close))
					return cache, nil
				},
			),
		)
	})

	return options
}
