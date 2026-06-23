package pubsub

import (
	"context"
	"sync"

	"cloud.google.com/go/pubsub/v2"
	v2 "github.com/cloudevents/sdk-go/v2"
	"gopkg.in/matryer/try.v1"

	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
	"github.com/xgodev/boost/wrapper/publisher"
)

// client implements a reusable Pub/Sub publisher.
type client struct {
	client  *pubsub.Client
	options *Options
	mu      sync.RWMutex
	topics  map[string]*pubsub.Publisher
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
		topics:  make(map[string]*pubsub.Publisher),
	}
}

// Publish sends a batch of events to Pub/Sub with concurrent goroutines and retry.
func (p *client) Publish(ctx context.Context, events []*v2.Event) ([]publisher.PublishOutput, error) {
	logger := log.FromContext(ctx).WithTypeOf(p)
	logger.Debug("publishing to Pub/Sub")

	if len(events) == 0 {
		logger.Warn("no messages to publish")
		return nil, nil
	}

	return p.send(ctx, events)
}

// send publishes events concurrently with retry and aggregates results as they arrive.
func (p *client) send(ctx context.Context, events []*v2.Event) ([]publisher.PublishOutput, error) {
	var wg sync.WaitGroup
	resultCh := make(chan publisher.PublishOutput, len(events))

	for _, ev := range events {
		wg.Add(1)
		go func(ev *v2.Event) {
			defer wg.Done()

			base := context.WithoutCancel(ctx)
			pubCtx, cancel := context.WithTimeout(base, p.options.PublishTimeout)
			defer cancel()

			logger := log.WithField("subject", ev.Subject()).
				WithField("id", ev.ID())

			msg := buildMessage(ev, p.options.OrderingKey)
			topic := p.getTopic(ev.Subject())

			err := try.Do(func(attempt int) (bool, error) {
				logger.Tracef("publishing to topic %s, attempt %d", ev.Subject(), attempt)
				r := topic.Publish(pubCtx, msg)
				if _, err := r.Get(pubCtx); err != nil {
					return attempt < 5, errors.NewInternal(err, "Pub/Sub publish failed")
				}
				logger.Infof("message published")
				return false, nil
			})

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
func (p *client) getTopic(subject string) *pubsub.Publisher {
	p.mu.RLock()
	if t, ok := p.topics[subject]; ok {
		p.mu.RUnlock()
		return t
	}
	p.mu.RUnlock()

	p.mu.Lock()
	defer p.mu.Unlock()

	if t, ok := p.topics[subject]; ok {
		return t
	}
	t := p.client.Publisher(subject)
	t.PublishSettings = p.options.Settings
	t.EnableMessageOrdering = p.options.OrderingKey

	p.topics[subject] = t
	return t
}

// Close stops all cached topics' background goroutines.
func (p *client) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, t := range p.topics {
		t.Stop()
	}
	p.topics = nil
}

type Output struct {
	publisher.PublishOutput
	Result *pubsub.PublishResult
}
