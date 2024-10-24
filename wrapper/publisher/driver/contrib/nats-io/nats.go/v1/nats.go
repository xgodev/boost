package nats

import (
	"context"
	"encoding/json"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
	"github.com/xgodev/boost/wrapper/publisher"
	"github.com/xgodev/boost/wrapper/publisher/util"

	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/nats-io/nats.go"
)

// client represents a NATS client that implements.
type client struct {
	conn *nats.Conn
}

// New creates a new NATS client.
func New(conn *nats.Conn) publisher.Driver {
	return &client{conn: conn}
}

// Publish publishes an event slice.
func (p *client) Publish(ctx context.Context, outs []*v2.Event) (res []publisher.PublishOutput, err error) {

	logger := log.FromContext(ctx).WithTypeOf(*p)

	logger.Tracef("publishing to nats")

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
					res = append(res, publisher.PublishOutput{Event: out, Error: errors.Wrap(err, errors.Internalf("unable to convert data to interface"))})
					continue
				}

				rawMessage, err = json.Marshal(data)

			} else {
				rawMessage, err = util.JSONBytes(*out)
			}

		} else {
			rawMessage, err = util.JSONBytes(*out)
		}

		if err != nil {
			res = append(res, publisher.PublishOutput{Event: out, Error: errors.Wrap(err, errors.Internalf("error when transforming json into bytes"))})
			continue
		}

		logger.Info(string(rawMessage))

		msg := &nats.Msg{
			Subject: out.Subject(),
			Data:    rawMessage,
		}

		err = p.conn.PublishMsg(msg)
		if err != nil {
			res = append(res, publisher.PublishOutput{Event: out, Error: errors.Wrap(err, errors.Internalf("unable to publish to nats"))})
			continue
		}

		res = append(res, publisher.PublishOutput{Event: out})
	}

	return res, err
}
