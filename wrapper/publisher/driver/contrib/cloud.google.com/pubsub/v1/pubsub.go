package pubsub

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/matryer/try"

	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
	"github.com/xgodev/boost/wrapper/publisher"
)

// client implements a reusable Pub/Sub publisher.
type client struct {
	client  *pubsub.Client
	options *Options

	mu     sync.Mutex
	topics map[string]*pubsub.Topic
}

// NewWithOptions returns a publisher with a topic cache.
func NewWithOptions(ctx context.Context, c *pubsub.Client, options *Options) publisher.Driver {
	return &client{
		client:  c,
		options: options,
		topics:  make(map[string]*pubsub.Topic),
	}
}

// Publish sends a batch of events to Pub/Sub.
func (p *client) Publish(ctx context.Context, events []*v2.Event) ([]publisher.PublishOutput, error) {
	logger := log.FromContext(ctx).WithTypeOf(*p)
	logger.Info("publishing to Pub/Sub")

	if len(events) == 0 {
		logger.Warn("no messages to publish")
		return nil, nil
	}
	return p.send(ctx, events)
}

// send iterates over the events and publishes them using cached topics.
func (p *client) send(ctx context.Context, events []*v2.Event) ([]publisher.PublishOutput, error) {
	var results []publisher.PublishOutput

	for _, ev := range events {
		logger := log.FromContext(ctx).WithTypeOf(*p).
			WithField("subject", ev.Subject()).
			WithField("id", ev.ID())

		// Convert event data to a generic map.
		var data map[string]interface{}
		if err := ev.DataAs(&data); err != nil {
			results = append(results, publisher.PublishOutput{
				Event: ev,
				Error: errors.Wrap(err, errors.Internalf("failed to convert event data")),
			})
			continue
		}

		// Serialize data to JSON.
		raw, err := json.Marshal(data)
		if err != nil {
			results = append(results, publisher.PublishOutput{
				Event: ev,
				Error: errors.Wrap(err, errors.Internalf("failed to marshal data")),
			})
			continue
		}

		// Build CloudEvents attributes.
		attrs := map[string]string{
			"ce_specversion": ev.SpecVersion(),
			"ce_id":          ev.ID(),
			"ce_source":      ev.Source(),
			"ce_type":        ev.Type(),
			"content-type":   ev.DataContentType(),
			"ce_time":        ev.Time().String(),
			"ce_path":        "/",
			"ce_subject":     ev.Subject(),
		}

		// Create Pub/Sub message.
		msg := &pubsub.Message{
			ID:          ev.ID(),
			Data:        raw,
			Attributes:  attrs,
			PublishTime: time.Now(),
		}

		// Set ordering key if enabled.
		if p.options.OrderingKey {
			pk, err := p.getPartitionKey(ev)
			if err != nil {
				results = append(results, publisher.PublishOutput{
					Event: ev,
					Error: errors.Wrap(err, errors.Internalf("failed to get partition key")),
				})
				continue
			}
			msg.OrderingKey = pk
		}

		// Retrieve cached topic or create a new one.
		topic := p.getTopic(ev.Subject())

		// Publish with retry logic.
		err = try.Do(func(attempt int) (bool, error) {
			logger.Tracef("publishing to topic %s, attempt %d", ev.Subject(), attempt)
			res := topic.Publish(ctx, msg)
			if _, err := res.Get(ctx); err != nil {
				log.Error(err)
				return attempt < 5, errors.NewInternal(err, "Pub/Sub publish failed")
			}
			logger.Infof("message published")
			logger.Debugf("payload: %s", string(raw))
			return false, nil
		})

		if err != nil {
			results = append(results, publisher.PublishOutput{Event: ev, Error: err})
		} else {
			results = append(results, publisher.PublishOutput{Event: ev})
		}
	}

	return results, nil
}

// getTopic returns a cached Pub/Sub topic or creates it on first use.
func (p *client) getTopic(subject string) *pubsub.Topic {
	p.mu.Lock()
	defer p.mu.Unlock()

	if t, ok := p.topics[subject]; ok {
		return t
	}
	t := p.client.Topic(subject)
	p.topics[subject] = t
	return t
}

// Close stops all cached topics' background goroutines.
// Call this when the publisher is shutting down.
func (p *client) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, t := range p.topics {
		t.Stop()
	}
	p.topics = nil
}

// getPartitionKey extracts the ordering key extension or uses the event ID.
func (p *client) getPartitionKey(ev *v2.Event) (string, error) {
	if key, ok := ev.Extensions()["key"]; ok {
		return key.(string), nil
	}
	return ev.ID(), nil
}
