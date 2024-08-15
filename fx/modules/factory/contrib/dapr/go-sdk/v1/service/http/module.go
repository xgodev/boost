package http

import (
	"github.com/xgodev/boost/factory/contrib/dapr/go-sdk/v1/service/http"
	"sync"

	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"go.uber.org/fx"
)

var once sync.Once

// Module fx module for dapr http service.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				http.New,
			),
		)
	})

	return options
}
