package nats

import (
	ginatsfx "github.com/xgodev/boost/factory/contrib/go.uber.org/fx/v1/module/contrib/nats-io/nats.go/v1"
	"sync"

	"github.com/xgodev/boost/factory/contrib/go.uber.org/fx/v1/module/core/context"
	"go.uber.org/fx"
)

var once sync.Once

// Module loads the NATS module providing an initialized client.
//
// The module is only loaded once.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			context.Module(),
			ginatsfx.Module(),
			fx.Provide(
				NewEvent,
			),
		)
	})

	return options
}
