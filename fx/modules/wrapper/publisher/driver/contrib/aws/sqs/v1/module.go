package sqs

import (
	"sync"

	clientfx "github.com/xgodev/boost/fx/modules/factory/contrib/aws/aws-sdk-go-v2/v1/client/sqs"
	sqs "github.com/xgodev/boost/wrapper/publisher/driver/contrib/aws/sqs/v1"
	"go.uber.org/fx"
)

var once sync.Once

// Module provides the sqs publisher.Driver on the boost sqs client.
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
			clientfx.Module(),
			fx.Provide(
				sqs.New,
			),
		)
	})

	return options
}
