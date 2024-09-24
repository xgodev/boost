package publisher

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
				publisher.Module(),
				fx.Provide(
					fx.Annotated{
						Group:  function.BSFunctionMiddlewaresGroupKey,
						Target: p.NewAnyErrorMiddleware[T],
					},
				),
			)
		})
	
	*/

	return options
}
