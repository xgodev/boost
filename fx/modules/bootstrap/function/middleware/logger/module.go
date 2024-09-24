package logger

import (
	"go.uber.org/fx"
	"sync"
)

var once sync.Once

func Module[T any]() fx.Option {
	options := fx.Options()
	if !IsEnabled() {
		return options
	}

	/*
		once.Do(func() {
			options = fx.Options(
				fx.Provide(logger.NewOptions),
				fx.Provide(
					fx.Annotated{
						Group:  function.BSFunctionMiddlewaresGroupKey,
						Target: logger.NewAnyErrorMiddlewareWithOptions[T],
					},
				),
			)
		})

	*/

	return options
}
