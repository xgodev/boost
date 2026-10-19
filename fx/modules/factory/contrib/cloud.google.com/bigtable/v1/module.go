package bigtable

import (
	"context"
	"sync"

	"cloud.google.com/go/bigtable"
	f "github.com/xgodev/boost/factory/contrib/cloud.google.com/bigtable/v1"
	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/client"
	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"go.uber.org/fx"
)

var once sync.Once

type params struct {
	fx.In
	Plugins []client.Plugin `optional:"true"`
}

// Module fx module for the Bigtable client of boost.factory.bigtable. The
// client is closed on app stop.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				func(ctx context.Context, lc fx.Lifecycle, p params) (*bigtable.Client, error) {
					c, err := f.NewClient(ctx, p.Plugins...)
					if err != nil {
						return nil, err
					}
					lc.Append(fx.StopHook(c.Close))
					return c, nil
				},
			),
		)
	})

	return options
}
