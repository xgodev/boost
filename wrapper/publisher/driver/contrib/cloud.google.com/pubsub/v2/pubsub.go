package pubsub

import (
	"context"
	"sync"

	pb "cloud.google.com/go/pubsub/v2"
	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/matryer/try"

	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
	"github.com/xgodev/boost/wrapper/publisher"
)

// client implements a reusable Pub/Sub publisher.
type client struct {
	client  *pb.Client
	options *Options

	mu     sync.Mutex
	topics map[string]*pb.Publisher
}

// NewWithConfigPath returns a publisher configured by a file path.
func NewWithConfigPath(ctx context.Context, c *pb.Client, path string) (publisher.Driver, error) {
	options, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewWithOptions(ctx, c, options), nil
}

// New returns a publisher with default options loaded from environment.
func New(ctx context.Context, c *pb.Client) (publisher.Driver, error) {
	options, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewWithOptions(ctx, c, options), nil
}

// NewWithOptions returns a publisher with explicit options.
func NewWithOptions(ctx context.Context, c *pb.Client, options *Options) publisher.Driver {
	return &client{
		client:  c,
		options: options,
		topics:  make(map[string]*pb.Publisher),
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

	for _, ev := range events {

		wg.Go(func() {
			pubsubMsg, err := generateCloudEvent(ev)
			if err != nil {
				resultCh <- publisher.PublishOutput{Event: ev, Error: err}
			}

			err = try.Do(func(attempt int) (bool, error) {
				log.Ctx(ctx, *p).WithField("subject", ev.Subject()).WithField("id", ev.ID()).Tracef("publishing to topic %s, attempt %d", ev.Subject(), attempt)
				topic := p.getPublisher(ev.Subject())
				r := topic.Publish(ctx, pubsubMsg)
				if _, err := r.Get(ctx); err != nil {
					log.Error(err)
					return attempt < 5, errors.NewInternal(err, "Pub/Sub publish failed")
				}

				log.Ctx(ctx, *p).WithField("subject", ev.Subject()).WithField("id", ev.ID()).Info("message published")
				return false, nil
			})

			resultCh <- publisher.PublishOutput{Event: ev, Error: err}
		})
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

	// Provavelmente deveria retornar erro aqui caso tenha algum erro na lista
	return results, nil
}

func (p *client) getPublisher(subject string) *pb.Publisher {
	p.mu.Lock()
	defer p.mu.Unlock()

	if t, ok := p.topics[subject]; ok {
		return t
	}
	t := p.client.Publisher(subject)
	t.PublishSettings = p.options.Settings
	p.topics[subject] = t
	return t
}

func (p *client) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, t := range p.topics {
		t.Stop()
	}
	p.topics = nil
}

func (p *client) getPartitionKey(ev *v2.Event) (string, error) {
	if key, ok := ev.Extensions()["key"]; ok {
		return key.(string), nil
	}
	return ev.ID(), nil
}
