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

// NewWithConfigPath returns a publisher configured by a file path.
func NewWithConfigPath(ctx context.Context, c *pubsub.Client, path string) (publisher.Driver, error) {
	options, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewWithOptions(ctx, c, options), nil
}

// New returns a publisher with default options loaded from environment.
func New(ctx context.Context, c *pubsub.Client) (publisher.Driver, error) {
	options, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewWithOptions(ctx, c, options), nil
}

// NewWithOptions returns a publisher with explicit options.
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

// send publishes events concurrently and aggregates results as they arrive.
func (p *client) send(ctx context.Context, events []*v2.Event) ([]publisher.PublishOutput, error) {
	var wg sync.WaitGroup
	resultCh := make(chan publisher.PublishOutput, len(events))

	// launch all publishes asynchronously
	for _, ev := range events {
		wg.Add(1)
		go func(ev *v2.Event) {
			defer wg.Done()

			logger := log.FromContext(ctx).WithTypeOf(*p).
				WithField("subject", ev.Subject()).
				WithField("id", ev.ID())

			// Convert event data
			var data map[string]interface{}
			if err := ev.DataAs(&data); err != nil {
				resultCh <- publisher.PublishOutput{Event: ev, Error: errors.Wrap(err, errors.Internalf("failed to convert event data"))}
				return
			}

			// Serialize to JSON
			raw, err := json.Marshal(data)
			if err != nil {
				resultCh <- publisher.PublishOutput{Event: ev, Error: errors.Wrap(err, errors.Internalf("failed to marshal data"))}
				return
			}

			// Build attributes
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

			msg := &pubsub.Message{ID: ev.ID(), Data: raw, Attributes: attrs, PublishTime: time.Now()}
			if p.options.OrderingKey {
				if pk, err := p.getPartitionKey(ev); err == nil {
					msg.OrderingKey = pk
				}
			}

			topic := p.getTopic(ev.Subject())
			err = try.Do(func(attempt int) (bool, error) {
				logger.Tracef("publishing to topic %s, attempt %d", ev.Subject(), attempt)
				r := topic.Publish(ctx, msg)
				if _, err := r.Get(ctx); err != nil {
					log.Error(err)
					return attempt < 5, errors.NewInternal(err, "Pub/Sub publish failed")
				}
				logger.Infof("message published")
				return false, nil
			})

			// send result as soon as done
			resultCh <- publisher.PublishOutput{Event: ev, Error: err}
		}(ev)
	}

	// close channel once all goroutines finish
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// collect results as they arrive
	var results []publisher.PublishOutput
	for res := range resultCh {
		results = append(results, res)
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
