package goka

import (
	"context"
	"encoding/json"
	"github.com/lovoo/goka"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
	"github.com/xgodev/boost/wrapper/publisher"
	"github.com/xgodev/boost/wrapper/publisher/util"

	v2 "github.com/cloudevents/sdk-go/v2"
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

		exts := out.Extensions()

		source, ok := exts["target"]

		if ok {

			s := source.(string)

			if s == "data" {
				var data interface{}

				err = out.DataAs(&data)
				if err != nil {
					return errors.Wrap(err, errors.Internalf("error on data as. %s", err.Error()))
				}

				rawMessage, err = json.Marshal(data)

			} else {
				rawMessage, err = util.JSONBytes(*out)
			}

		} else {
			rawMessage, err = util.JSONBytes(*out)
		}

		if err != nil {
			return errors.Wrap(err, errors.Internalf("error when transforming json into bytes"))
		}

		pk, err := p.partitionKey(out)
		if err != nil {
			return errors.Wrap(err, errors.Internalf("unable to gets partition key"))
		}

		err = p.emitter.EmitSync(pk, rawMessage)
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
