package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/matryer/try"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
	"golang.org/x/sync/errgroup"
)

// Client represents a sns client.
type Client struct {
	client *pubsub.Client
}

// NewClient creates a new sns client.
func NewClient(c *pubsub.Client) *Client {
	return &Client{client: c}
}

// Publish publishes an event slice.
func (p *Client) Publish(ctx context.Context, events []*v2.Event) error {

	logger := log.FromContext(ctx).WithTypeOf(*p)

	logger.Info("publishing to awssns")

	if len(events) > 0 {

		return p.send(ctx, events)

	}

	logger.Warnf("no messages were reported for posting")

	return nil
}

func (p *Client) send(parentCtx context.Context, events []*v2.Event) (err error) {

	logger := log.FromContext(parentCtx).WithTypeOf(*p)

	g, gctx := errgroup.WithContext(parentCtx)
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
					return attempt < 5, errors.NewInternal(err, "could not be published in gcp pubsub")
				}
				return false, nil
			})

			return err

		})

	}

	return g.Wait()
}
