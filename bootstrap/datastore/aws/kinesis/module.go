package kinesis

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/xgodev/boost/bootstrap/datastore/aws"
	igkinesis "github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/client/kinesis"
	"github.com/xgodev/boost/factory/contrib/go.uber.org/fx/v1/module/core/context"
	"go.uber.org/fx"
)

var once sync.Once

// Module loads the kinesis module providing an initialized client.
//
// The module is only loaded once.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			context.Module(),
			aws.Module(),
			fx.Provide(
				kinesis.NewFromConfig,
				igkinesis.NewClient,
				DefaultOptions,
				NewEvent,
			),
		)
	})

	return options
}
