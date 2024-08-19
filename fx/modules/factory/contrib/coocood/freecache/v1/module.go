package freecache

import (
	"github.com/xgodev/boost/factory/contrib/coocood/freecache/v1"
	"sync"

	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"go.uber.org/fx"
)

var once sync.Once

// Module fx module for freecache.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				freecache.NewCache,
			),
		)
	})

	return options
}
