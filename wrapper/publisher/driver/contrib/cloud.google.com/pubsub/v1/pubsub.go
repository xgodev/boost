package pubsub

import (
	"context"
	"encoding/json"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
	"github.com/xgodev/boost/wrapper/publisher"
	"time"

	"cloud.google.com/go/pubsub"
	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/matryer/try"
)

// client represents a pubsub client.
type client struct {
	client  *pubsub.Client
	options *Options
}

// NewWithConfigPath returns connection with options from config path.
func NewWithConfigPath(ctx context.Context, c *pubsub.Client, path string) (publisher.Driver, error) {
	options, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewWithOptions(ctx, c, options), nil
}

// NewWithOptions returns connection with options.
func NewWithOptions(ctx context.Context, c *pubsub.Client, options *Options) publisher.Driver {
	return &client{options: options, client: c}
}

// New creates a new pubsub client.
func New(ctx context.Context, c *pubsub.Client) (publisher.Driver, error) {

	options, err := NewOptions()
	if err != nil {
		return nil, err
	}

	return NewWithOptions(ctx, c, options), nil
}

// Publish publishes an event slice.
func (p *client) Publish(ctx context.Context, events []*v2.Event) ([]publisher.PublishOutput, error) {

	logger := log.FromContext(ctx).WithTypeOf(*p)

	logger.Info("publishing to pubsub")

	if len(events) > 0 {

		return p.send(ctx, events)

	}

	logger.Warnf("no messages were reported for posting")

	return []publisher.PublishOutput{}, nil
}

func (p *client) send(ctx context.Context, events []*v2.Event) (res []publisher.PublishOutput, err error) {

	logger := log.FromContext(ctx).WithTypeOf(*p)

	for _, out := range events {

		var data map[string]interface{}
		if err := out.DataAs(&data); err != nil {
			res = append(res, publisher.PublishOutput{Event: out, Error: errors.Wrap(err, errors.Internalf("unable to convert data to interface"))})
			continue
		}

		var rawMessage []byte
		rawMessage, err = json.Marshal(data)
		if err != nil {
			res = append(res, publisher.PublishOutput{Event: out, Error: errors.Wrap(err, errors.Internalf("error on marshal"))})
			continue
		}

		attrs := map[string]string{
			"ce_specversion": out.SpecVersion(),
			"ce_id":          out.ID(),
			"ce_source":      out.Source(),
			"ce_type":        out.Type(),
			"content-type":   out.DataContentType(),
			"ce_time":        out.Time().String(),
			"ce_path":        "/",
			"ce_subject":     out.Subject(),
		}

		message := &pubsub.Message{
			ID:              out.ID(),
			Data:            rawMessage,
			Attributes:      attrs,
			PublishTime:     time.Now(),
			DeliveryAttempt: nil,
		}

		if p.options.OrderingKey {
			pk, err := p.partitionKey(out)
			if err != nil {
				res = append(res, publisher.PublishOutput{Event: out, Error: errors.Wrap(err, errors.Internalf("unable to gets partition key"))})
				continue
			}
			message.OrderingKey = pk
		}

		topic := p.client.Topic(out.Subject())
		defer topic.Stop()

		l := logger.WithField("subject", out.Subject()).
			WithField("id", out.ID())

		err = try.Do(func(attempt int) (bool, error) {

			l.Tracef("publishing message to topic %s attempt %v", out.Subject(), attempt)

			r := topic.Publish(ctx, message)
			if _, err := r.Get(ctx); err != nil {
				log.Error(err)
				return attempt < 5, errors.NewInternal(err, "could not be published in gcp pubsub")
			}
			l.Infof("message published")
			l.Debugf(string(rawMessage))
			return false, nil
		})

		if err != nil {
			res = append(res, publisher.PublishOutput{Event: out, Error: err})
			continue
		}

		res = append(res, publisher.PublishOutput{Event: out})

	}

	return res, err
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
