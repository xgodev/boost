package aws

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	a "github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1"
	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"go.uber.org/fx"
)

var once sync.Once

type params struct {
	fx.In
	Plugins []a.Plugin `optional:"true"`
}

// Module fx module for aws config.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				func(ctx context.Context, p params) (aws.Config, error) {
					return a.NewConfig(ctx, p.Plugins...)
				},
			),
		)
	})

	return options
}
