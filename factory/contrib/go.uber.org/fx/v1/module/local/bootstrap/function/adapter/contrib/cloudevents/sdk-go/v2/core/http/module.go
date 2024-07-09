package http

import (
	ce "github.com/cloudevents/sdk-go/v2"
	fn "github.com/xgodev/boost/bootstrap/function"
	"github.com/xgodev/boost/bootstrap/function/adapter/contrib/cloudevents/sdk-go/v2/core/http"
	fxnats "github.com/xgodev/boost/factory/contrib/go.uber.org/fx/v1/module/contrib/nats-io/nats.go/v1"
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
			fxnats.Module(),
			fx.Provide(
				fx.Annotated{
					Group: function.BSFunctionAdaptersGroupKey,
					Target: func() fn.CmdFunc {
						return http.New(
							ce.WithUUIDs(),
							ce.WithTimeNow(),
						)
					},
				},
			),
		)
	})

	return options
}
