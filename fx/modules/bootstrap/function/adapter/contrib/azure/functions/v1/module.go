package azure

import (
	"sync"

	"github.com/xgodev/boost/bootstrap/function/adapter/contrib/azure/functions/v1"
	"github.com/xgodev/boost/fx/modules/bootstrap/function"
	"go.uber.org/fx"
)

var once sync.Once

// Module adds the adapter that runs the function as an Azure Functions custom
// handler (the azure command) to the function adapters.
func Module[T any]() fx.Option {
	options := fx.Options()
	if !IsEnabled() {
		return options
	}

	once.Do(func() {
		options = fx.Options(
			fx.Provide(
				fx.Annotated{
					Group:  function.BSFunctionAdaptersGroupKey,
					Target: azure.New[T],
				},
			),
		)
	})

	return options
}
