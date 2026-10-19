package sns

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	client "github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/client/sns"
	awsfx "github.com/xgodev/boost/fx/modules/factory/contrib/aws/aws-sdk-go-v2/v1"
	"go.uber.org/fx"
)

var once sync.Once

// Module provides the SNS client of the SDK (*sns.Client), built from the
// aws.Config of the aws module, and the boost client.Client on it.
//
// It returns an empty option after the first call.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			awsfx.Module(),
			fx.Provide(
				func(cfg aws.Config) *sns.Client { return sns.NewFromConfig(cfg) },
				client.NewClient,
			),
		)
	})

	return options
}
