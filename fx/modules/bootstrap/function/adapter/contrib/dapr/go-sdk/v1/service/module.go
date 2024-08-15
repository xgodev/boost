package service

import (
	"github.com/xgodev/boost/bootstrap/function/adapter/contrib/dapr/go-sdk/v1/service"
	"github.com/xgodev/boost/fx/modules/bootstrap/function"
	"github.com/xgodev/boost/fx/modules/factory/contrib/dapr/go-sdk/v1/service/grpc"
	"github.com/xgodev/boost/fx/modules/factory/contrib/dapr/go-sdk/v1/service/http"
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

		var tp fx.Option
		switch Type() {
		case "grpc":
			tp = grpc.Module()
		default:
			tp = http.Module()
		}

		options = fx.Options(
			tp,
			fx.Provide(
				fx.Annotated{
					Group:  function.BSFunctionAdaptersGroupKey,
					Target: service.New[T],
				},
			),
		)
	})

	return options
}
