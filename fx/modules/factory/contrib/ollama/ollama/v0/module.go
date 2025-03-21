package ollama

import (
	"github.com/xgodev/boost/factory/contrib/ollama/ollama/v0"
	"sync"

	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"go.uber.org/fx"
)

var once sync.Once

// Module fx module for ollama client.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {

		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				ollama.NewClient,
			),
		)
	})

	return options
}
