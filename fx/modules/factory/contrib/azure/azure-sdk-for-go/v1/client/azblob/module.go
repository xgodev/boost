package azblob

import (
	"context"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	client "github.com/xgodev/boost/factory/contrib/azure/azure-sdk-for-go/v1/client/azblob"
	azurefx "github.com/xgodev/boost/fx/modules/factory/contrib/azure/azure-sdk-for-go/v1"
	"go.uber.org/fx"
)

type params struct {
	fx.In
	Credential azcore.TokenCredential
	Plugins    []client.Plugin `optional:"true"`
}

var once sync.Once

// Module fx module for the azblob client. It provides *azblob.Client, built
// with the azure credential module's credential.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			azurefx.Module(),
			fx.Provide(
				func(ctx context.Context, p params) (*azblob.Client, error) {
					return client.NewClient(ctx, p.Credential, p.Plugins...)
				},
			),
		)
	})

	return options
}
