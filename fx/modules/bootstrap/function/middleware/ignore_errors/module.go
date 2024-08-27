package ignore_errors

import (
	p "github.com/xgodev/boost/bootstrap/function/middleware/ignore_errors"
	"github.com/xgodev/boost/fx/modules/bootstrap/function"
	"go.uber.org/fx"
	"sync"
)

var once sync.Once

func Module[T any]() fx.Option {
	options := fx.Options()
	if !IsEnabled() {
		return options
	}

	once.Do(func() {
		options = fx.Options(
			fx.Provide(
				fx.Annotated{
					Group:  function.BSFunctionMiddlewaresGroupKey,
					Target: p.NewAnyErrorMiddleware[T],
				},
			),
		)
	})

	return options
}
