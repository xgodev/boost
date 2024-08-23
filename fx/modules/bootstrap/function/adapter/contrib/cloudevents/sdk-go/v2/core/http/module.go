package http

import (
	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
	fn "github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/bootstrap/function/adapter/contrib/cloudevents/sdk-go/v2/core/http"
	"github.com/xgodev/boost/fx/modules/bootstrap/function"
	"go.uber.org/fx"
	"sync"
)

type params struct {
	fx.In
	Plugins []http.Plugin `optional:"true"`
}

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
					Group: function.BSFunctionAdaptersGroupKey,
					Target: func(p params) fn.CmdFunc[T] {
						return http.New[T](
							[]client.Option{
								ce.WithUUIDs(),
								ce.WithTimeNow(),
							},
							p.Plugins...,
						)
					},
				},
			),
		)
	})

	return options
}
