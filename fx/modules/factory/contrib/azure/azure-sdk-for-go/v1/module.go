package azure

import (
	"context"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azure "github.com/xgodev/boost/factory/contrib/azure/azure-sdk-for-go/v1"
	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"go.uber.org/fx"
)

var once sync.Once

// Module fx module for the azure credential. It provides azcore.TokenCredential.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				func(ctx context.Context) (azcore.TokenCredential, error) {
					return azure.NewCredential(ctx)
				},
			),
		)
	})

	return options
}
