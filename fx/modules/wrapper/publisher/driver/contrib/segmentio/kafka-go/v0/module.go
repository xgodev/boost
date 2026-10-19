package kafka

import (
	"sync"

	factoryfx "github.com/xgodev/boost/fx/modules/factory/contrib/segmentio/kafka-go/v0"
	kafka "github.com/xgodev/boost/wrapper/publisher/driver/contrib/segmentio/kafka-go/v0"
	"go.uber.org/fx"
)

var once sync.Once

// Module provides the kafka-go publisher.Driver on the boost kafka-go writer.
//
// It returns an empty option when the driver is disabled or after the first
// call.
func Module() fx.Option {
	options := fx.Options()
	if !IsEnabled() {
		return options
	}

	once.Do(func() {
		options = fx.Options(
			factoryfx.WriterModule(),
			fx.Provide(
				kafka.New,
			),
		)
	})

	return options
}
