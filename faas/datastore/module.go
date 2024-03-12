package datastore

import (
	"github.com/xgodev/boost/faas/datastore/gcp/pubsub"
	"sync"

	"github.com/xgodev/boost/faas/datastore/aws/kinesis"
	"github.com/xgodev/boost/faas/datastore/aws/sns"
	"github.com/xgodev/boost/faas/datastore/aws/sqs"
	"github.com/xgodev/boost/faas/datastore/nats"
	"github.com/xgodev/boost/faas/repository"
	"github.com/xgodev/boost/log"
	"go.uber.org/fx"
)

var eventOnce sync.Once

// EventModule returns fx module for initialization of event module based on the configured event provider.
// It can be: nats (default), kinesis, sns or sqs.
//
// The module is only loaded once.
func EventModule() fx.Option {

	options := fx.Options()

	eventOnce.Do(func() {

		value := repository.EventProviderValue()
		log.Debugf("Loading event provider %s", value)
		switch value {
		case "kinesis":
			options = kinesis.Module()
		case "sns":
			options = sns.Module()
		case "sqs":
			options = sqs.Module()
		case "pubsub":
			options = pubsub.Module()
		default:
			options = nats.Module()
		}

	})

	return options
}
