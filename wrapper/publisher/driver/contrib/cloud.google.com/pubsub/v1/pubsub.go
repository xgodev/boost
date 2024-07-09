package pubsub

import (
	"context"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
	"github.com/xgodev/boost/wrapper/publisher"

	"cloud.google.com/go/pubsub"
	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/matryer/try"
	"golang.org/x/sync/errgroup"
)

// client represents a pubsub client.
type client struct {
	client *pubsub.Client
}

// New creates a new pubsub client.
func New(c *pubsub.Client) publisher.Driver {
	return &client{client: c}
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

		event := e

		g.Go(func() (err error) {

			var rawMessage []byte

			rawMessage, err = event.MarshalJSON()
			if err != nil {
				return errors.Wrap(err, errors.Internalf("error on marshal. %s", err.Error()))
			}

			message := &pubsub.Message{
				Data: rawMessage,
			}

			topic := p.client.Topic(event.Subject())
			defer topic.Stop()

			logger.WithField("subject", event.Subject()).
				WithField("id", event.ID()).
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
