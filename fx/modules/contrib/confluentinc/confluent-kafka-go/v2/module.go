package confluent

import (
	"github.com/xgodev/boost/factory/contrib/confluentinc/confluent-kafka-go/v2"
	"sync"

	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"go.uber.org/fx"
)

var once sync.Once

// ProducerModule fx module for kafka connection.
func ProducerModule() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			contextfx.Module(),
			fx.Provide(
				confluent.NewProducer,
			),
		)
	})

	return options
}
