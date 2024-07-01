package datastore

import (
	"github.com/xgodev/boost/bootstrap/datastore/gcp/pubsub"
	"sync"

	"github.com/xgodev/boost/bootstrap/datastore/aws/kinesis"
	"github.com/xgodev/boost/bootstrap/datastore/aws/sns"
	"github.com/xgodev/boost/bootstrap/datastore/aws/sqs"
	"github.com/xgodev/boost/bootstrap/datastore/nats"
	"github.com/xgodev/boost/bootstrap/repository"
	"github.com/xgodev/boost/wrapper/log"
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
