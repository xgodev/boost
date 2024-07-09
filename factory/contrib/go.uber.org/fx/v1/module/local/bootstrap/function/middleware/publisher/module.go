package publisher

import (
	"github.com/xgodev/boost/bootstrap/function/middleware/publisher"
	"github.com/xgodev/boost/factory/contrib/go.uber.org/fx/v1/module/local/bootstrap/function"
	"go.uber.org/fx"
	"sync"
)

var once sync.Once

func Module() fx.Option {
	options := fx.Options()
	if !IsEnabled() {
		return options
	}

	once.Do(func() {
		options = fx.Options(
			fx.Provide(
				fx.Annotated{
					Group:  function.BSFunctionMiddlewaresGroupKey,
					Target: publisher.New,
				},
			),
		)
	})

	return options
}
