package lambda

import (
	"sync"

	"github.com/xgodev/boost/bootstrap/function/adapter/contrib/aws/aws-lambda-go/v1"
	"github.com/xgodev/boost/fx/modules/bootstrap/function"
	"go.uber.org/fx"
)

var once sync.Once

// Module adds the adapter that runs the function as an AWS Lambda function
// (the lambda command) to the function adapters.
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
					Target: lambda.New[T],
				},
			),
		)
	})

	return options
}
