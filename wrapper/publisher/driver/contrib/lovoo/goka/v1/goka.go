package goka

import (
	"context"
	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/xgodev/boost/factory/contrib/lovoo/goka/v1"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
	"github.com/xgodev/boost/wrapper/publisher"
)

// client represents a Kafka client that implements.
type client struct {
	emitter *goka.Emitter
}

// New creates a new Kafka client.
func New(emitter *goka.Emitter) publisher.Driver {
	return &client{emitter: emitter}
}

// Publish publishes an event slice.
func (p *client) Publish(ctx context.Context, outs []*v2.Event) (err error) {

	logger := log.FromContext(ctx).WithTypeOf(*p)

	logger.Tracef("publishing to kafka")

	for _, out := range outs {

		logger = logger.
			WithField("subject", out.Subject()).
			WithField("id", out.ID())

		var rawMessage []byte

		rawMessage, err = out.MarshalJSON()
		if err != nil {
			return errors.Wrap(err, errors.Internalf("error on marshal. %s", err.Error()))
		}

		pk, err := p.partitionKey(out)
		if err != nil {
			return errors.Wrap(err, errors.Internalf("unable to gets partition key"))
		}

		err = p.emitter.EmitSync(ctx, out.Subject(), pk, rawMessage)
		if err != nil {
			return errors.Wrap(err, errors.Internalf("unable to publish to kafka"))
		}

		logger.Info(string(rawMessage))

	}

	return nil
}

func (p *client) partitionKey(out *v2.Event) (string, error) {

	var pk string
	exts := out.Extensions()

	if key, ok := exts["key"]; ok {
		pk = key.(string)
	} else {
		pk = out.ID()
	}

	return pk, nil
}
