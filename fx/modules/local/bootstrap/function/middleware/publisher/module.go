package publisher

import (
	p "github.com/xgodev/boost/bootstrap/function/middleware/publisher"
	"github.com/xgodev/boost/fx/modules/local/bootstrap/function"
	"github.com/xgodev/boost/fx/modules/local/extra/publisher"
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
			publisher.Module(),
			fx.Provide(
				fx.Annotated{
					Group:  function.BSFunctionMiddlewaresGroupKey,
					Target: p.New,
				},
			),
		)
	})

	return options
}
