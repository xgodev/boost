package sns

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/xgodev/boost/bootstrap/datastore/aws"
	igsns "github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/client/sns"
	"github.com/xgodev/boost/factory/contrib/go.uber.org/fx/v1/module/core/context"
	"go.uber.org/fx"
)

var once sync.Once

// Module loads the sns module providing an initialized client.
//
// The module is only loaded once.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			context.Module(),
			aws.Module(),
			fx.Provide(
				sns.NewFromConfig,
				igsns.NewClient,
				NewEvent,
			),
		)
	})

	return options
}
