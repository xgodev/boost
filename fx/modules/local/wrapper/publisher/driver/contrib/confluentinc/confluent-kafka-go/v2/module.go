package confluent

import (
	fxconfluent "github.com/xgodev/boost/fx/modules/contrib/confluentinc/confluent-kafka-go/v2"
	"github.com/xgodev/boost/wrapper/publisher/driver/contrib/confluentinc/confluent-kafka-go/v2"
	"go.uber.org/fx"
	"sync"
)

var once sync.Once

func Module() fx.Option {
	options := fx.Options()
	if !IsEnabled() {
		return options
	}

	once.Do(func() {
		options = fx.Options(
			fxconfluent.ProducerModule(),
			fx.Provide(
				confluent.New,
			),
		)
	})

	return options
}
