package pubsub

import (
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsub/pstest"
	"context"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"sync/atomic"
	"testing"
	"time"
)

func setup(t *testing.T) (*pubsub.Client, context.Context) {
	srv := pstest.NewServer()
	t.Setenv("PUBSUB_EMULATOR_HOST", srv.Addr)

	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, "local-client")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	return client, ctx
}

func TestConcurrencyInSubscriber(t *testing.T) {

	client, ctx := setup(t)

	err := setupTopicAndSubscription(t, client, ctx, 30)
	if err != nil {
		t.Fatalf("Failed to setup test: %v", err)
	}

	received := runSubscriber(t, ctx, client, 7, 4*time.Second)
	var expected int32 = 28
	if received != expected {
		t.Errorf("received %d messages, expected %d", received, expected)
	}

}

func runSubscriber(t *testing.T, ctx context.Context, client *pubsub.Client, concurrency int64, duration time.Duration) int32 {
	h := Handler{sleep: 1 * time.Second}
	sub := NewSubscriber[cloudevents.Event](client, h.handle, "subscription-test", &Options{Concurrency: concurrency})

	ctxTimeout, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

	err := sub.Subscribe(ctxTimeout)
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	return h.received
}

func setupTopicAndSubscription(t *testing.T, client *pubsub.Client, ctx context.Context, numMessages int) error {
	to, err := client.CreateTopic(ctx, "topic-test")
	if err != nil {
		t.Fatalf("Failed to create topic: %v", err)
	}

	_, err = client.CreateSubscription(ctx, "subscription-test", pubsub.SubscriptionConfig{
		Topic:                 to,
		AckDeadline:           10 * time.Second,
		EnableMessageOrdering: true,
	})
	if err != nil {
		t.Fatalf("Failed to create subscription: %v", err)
	}

	err = publishMsgs(ctx, to, numMessages)

	return err
}

type Handler struct {
	received int32
	sleep    time.Duration
}

func (h *Handler) handle(ctx context.Context, event cloudevents.Event) (cloudevents.Event, error) {
	atomic.AddInt32(&h.received, 1)
	time.Sleep(h.sleep)

	fmt.Println("processing message")
	return cloudevents.Event{}, nil
}

func publishMsgs(ctx context.Context, t *pubsub.Topic, numMsgs int) error {
	var results []*pubsub.PublishResult
	t.PublishSettings.CountThreshold = 1

	for i := 0; i < numMsgs; i++ {
		res := t.Publish(ctx, &pubsub.Message{
			Data: []byte(fmt.Sprintf("message#%d", i)),
		})
		results = append(results, res)
	}

	for _, r := range results {
		if _, err := r.Get(ctx); err != nil {
			return fmt.Errorf("Get publish result: %w", err)
		}
	}

	return nil
}
