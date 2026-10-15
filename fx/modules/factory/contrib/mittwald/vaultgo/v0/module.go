package vault

import (
	"sync"

	v "github.com/xgodev/boost/factory/contrib/mittwald/vaultgo/v0"
	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"go.uber.org/fx"
)

var once sync.Once

// Module fx module for vault client.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				v.NewClient,
			),
		)
	})

	return options
}
