package sqs

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/xgodev/boost/bootstrap/datastore/aws"
	igsqs "github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/client/sqs"
	"github.com/xgodev/boost/factory/contrib/go.uber.org/fx/v1/module/core/context"
	"go.uber.org/fx"
)

var once sync.Once

// Module loads the sqs module providing an initialized client.
//
// The module is only loaded once.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {
		options = fx.Options(
			context.Module(),
			aws.Module(),
			fx.Provide(
				sqs.NewFromConfig,
				igsqs.NewClient,
				NewEvent,
			),
		)
	})

	return options
}
