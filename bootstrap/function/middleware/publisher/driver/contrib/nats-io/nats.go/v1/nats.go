package nats

import (
	"context"
	"encoding/json"
	"github.com/xgodev/boost/bootstrap/function/middleware/publisher"
	"github.com/xgodev/boost/bootstrap/function/util"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"

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
func (p *client) Publish(ctx context.Context, outs []*v2.Event) (err error) {

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
			err = errors.Wrap(err, errors.Internalf("error when transforming json into bytes"))
			logger.Error(errors.ErrorStack(err))
			continue
		}

		logger.Info(string(rawMessage))

		msg := &nats.Msg{
			Subject: out.Subject(),
			Data:    rawMessage,
		}

		err = p.conn.PublishMsg(msg)
		if err != nil {
			err = errors.Wrap(err, errors.Internalf("unable to publish to nats"))
			logger.Error(errors.ErrorStack(err))
		}

	}

	return nil
}
