package kinesis

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	client "github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/client/kinesis"
	awsfx "github.com/xgodev/boost/fx/modules/factory/contrib/aws/aws-sdk-go-v2/v1"
	"go.uber.org/fx"
)

var once sync.Once

// Module provides the Kinesis client of the SDK (*kinesis.Client), built from the
// aws.Config of the aws module, and the boost client.Client on it.
//
// It returns an empty option after the first call.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			awsfx.Module(),
			fx.Provide(
				func(cfg aws.Config) *kinesis.Client { return kinesis.NewFromConfig(cfg) },
				client.NewClient,
			),
		)
	})

	return options
}
