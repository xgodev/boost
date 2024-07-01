package aws

import (
	"sync"

	igaws "github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1"
	"go.uber.org/fx"
)

var once sync.Once

// Module returns fx module for initialization of aws config.
//
// The module is only loaded once.
func Module() fx.Option {

	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			fx.Provide(
				igaws.NewConfig,
			),
		)
	})

	return options
}
