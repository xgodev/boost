package log

import (
	"github.com/xgodev/boost/bootstrap/function/middleware/log"
	"github.com/xgodev/boost/fx/modules/local/bootstrap/function"
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
			fx.Provide(log.NewOptions),
			fx.Provide(
				fx.Annotated{
					Group:  function.BSFunctionMiddlewaresGroupKey,
					Target: log.New,
				},
			),
		)
	})

	return options
}
