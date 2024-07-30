package grpc

import (
	"github.com/xgodev/boost/factory/contrib/dapr/go-sdk/v1/service/grpc"
	"sync"

	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"go.uber.org/fx"
)

var once sync.Once

// Module fx module for dapr grpc connection.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				grpc.NewService,
			),
		)
	})

	return options
}
