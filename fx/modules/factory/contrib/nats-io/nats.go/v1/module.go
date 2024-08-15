package nats

import (
	"github.com/xgodev/boost/factory/contrib/nats-io/nats.go/v1"
	"sync"

	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"go.uber.org/fx"
)

var once sync.Once

// Module fx module for nats.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				nats.NewConn,
			),
		)
	})

	return options
}
