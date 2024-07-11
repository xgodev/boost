package goka

import (
	"github.com/xgodev/boost/factory/contrib/lovoo/goka/v1"
	"sync"

	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"go.uber.org/fx"
)

var once sync.Once

// Module fx module for kafka connection.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				goka.NewEmitter,
			),
		)
	})

	return options
}
