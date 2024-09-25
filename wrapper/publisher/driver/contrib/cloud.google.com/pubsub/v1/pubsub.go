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
	"golang.org/x/sync/errgroup"
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
func (p *client) Publish(ctx context.Context, events []*v2.Event) error {

	logger := log.FromContext(ctx).WithTypeOf(*p)

	logger.Info("publishing to pubsub")

	if len(events) > 0 {

		return p.send(ctx, events)

	}

	logger.Warnf("no messages were reported for posting")

	return nil
}

func (p *client) send(ctx context.Context, events []*v2.Event) (err error) {

	logger := log.FromContext(ctx).WithTypeOf(*p)

	g, gctx := errgroup.WithContext(ctx)
	defer gctx.Done()

	for _, e := range events {

		out := e

		g.Go(func() (err error) {

			var data map[string]interface{}
			if err := out.DataAs(&data); err != nil {
				return errors.Wrap(err, errors.Internalf("error on marshal. %s", err.Error()))
			}

			var rawMessage []byte
			rawMessage, err = json.Marshal(data)
			if err != nil {
				return errors.Wrap(err, errors.Internalf("error on marshal. %s", err.Error()))
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

			pk, err := p.partitionKey(out)
			if err != nil {
				return errors.Wrap(err, errors.Internalf("unable to gets partition key"))
			}

			message := &pubsub.Message{
				ID:              out.ID(),
				Data:            rawMessage,
				Attributes:      attrs,
				PublishTime:     time.Now(),
				DeliveryAttempt: nil,
				OrderingKey:     pk,
			}

			topic := p.client.Topic(out.Subject())
			defer topic.Stop()

			logger.WithField("subject", out.Subject()).
				WithField("id", out.ID()).
				Info(string(rawMessage))

			err = try.Do(func(attempt int) (bool, error) {
				r := topic.Publish(gctx, message)
				if _, err := r.Get(gctx); err != nil {
					log.Error(err)
					return attempt < 5, errors.NewInternal(err, "could not be published in gcp pubsub")
				}
				return false, nil
			})

			return err

		})

	}

	return g.Wait()
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
