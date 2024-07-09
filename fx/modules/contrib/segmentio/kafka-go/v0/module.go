package kafka

import (
	"github.com/xgodev/boost/factory/contrib/segmentio/kafka-go/v0"
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
				kafka.NewConn,
			),
		)
	})

	return options
}
