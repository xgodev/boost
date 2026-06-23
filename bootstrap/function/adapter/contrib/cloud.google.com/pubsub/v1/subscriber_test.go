package pubsub

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"cloud.google.com/go/pubsub/pstest"
	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"google.golang.org/protobuf/types/known/durationpb"
)

func setup(t *testing.T) (*pubsub.Client, context.Context) {
	t.Helper()

	srv := pstest.NewServer()
	t.Setenv("PUBSUB_EMULATOR_HOST", srv.Addr)
	t.Cleanup(func() {
		if err := srv.Close(); err != nil {
			t.Logf("error closing pstest server: %v", err)
		}
	})

	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, "local-client")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	t.Cleanup(func() {
		if err := client.Close(); err != nil {
			t.Logf("error closing pubsub client: %v", err)
		}
	})

	return client, ctx
}

func TestConcurrencyInSubscriber(t *testing.T) {
	tests := []struct {
		name        string
		concurrency int64
		duration    time.Duration
		expected    int32
	}{
		{
			name:        "7 workers for 4 seconds",
			concurrency: 7,
			duration:    4 * time.Second,
			expected:    28, // 7 workers × 4s / 1s per msg
		},
		{
			name:        "single worker",
			concurrency: 1,
			duration:    3 * time.Second,
			expected:    3, // 1 worker × 3s / 1s per msg
		},
		{
			name:        "high concurrency",
			concurrency: 20,
			duration:    2 * time.Second,
			expected:    30, // 20 workers × 2 messages in 2s, capped at 30 messages total
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, ctx := setup(t)

			const numMessages = 30
			subName, err := setupTopicAndSubscription(client, ctx, numMessages, tt.name)
			if err != nil {
				t.Fatalf("Failed to setup test: %v", err)
			}

			received := runSubscriber(t, ctx, client, subName, tt.concurrency, tt.duration)
			if received != tt.expected {
				t.Errorf("received %d messages, expected %d", received, tt.expected)
			}
		})
	}
}

func runSubscriber(t *testing.T, ctx context.Context, client *pubsub.Client, subscription string, concurrency int64, duration time.Duration) int32 {
	t.Helper()

	h := Handler{sleep: 1 * time.Second}
	sub := NewSubscriber[cloudevents.Event](client, h.handle, subscription, &Options{Concurrency: concurrency})

	ctxTimeout, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

	err := sub.Subscribe(ctxTimeout)
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	return h.received
}

func setupTopicAndSubscription(client *pubsub.Client, ctx context.Context, numMessages int, suffix string) (string, error) {
	topicName := "topic-test-" + suffix
	subName := "subscription-test-" + suffix
	topicPath := "projects/local-client/topics/" + topicName
	subPath := "projects/local-client/subscriptions/" + subName

	_, err := client.TopicAdminClient.CreateTopic(ctx, &pubsubpb.Topic{
		Name: topicPath,
	})
	if err != nil {
		return subName, fmt.Errorf("create topic: %w", err)
	}

	_, err = client.SubscriptionAdminClient.CreateSubscription(ctx, &pubsubpb.Subscription{
		Name:                     subPath,
		Topic:                    topicPath,
		AckDeadlineSeconds:       10,
		EnableMessageOrdering:    true,
		MessageRetentionDuration: durationpb.New(10 * time.Minute),
	})
	if err != nil {
		return subName, fmt.Errorf("create subscription: %w", err)
	}

	return subName, publishMsgs(ctx, client, topicName, numMessages)
}

type Handler struct {
	received      int32
	sleep         time.Duration
	shouldSucceed func() bool
}

func (h *Handler) handle(_ context.Context, _ cloudevents.Event) (cloudevents.Event, error) {
	if h.shouldSucceed != nil && !h.shouldSucceed() {
		return cloudevents.Event{}, fmt.Errorf("handler error")
	}

	atomic.AddInt32(&h.received, 1)
	time.Sleep(h.sleep)

	return cloudevents.Event{}, nil
}

func publishMsgs(ctx context.Context, client *pubsub.Client, topicName string, numMsgs int) error {
	var results []*pubsub.PublishResult
	pub := client.Publisher(topicName)
	pub.PublishSettings.CountThreshold = 1

	for i := 0; i < numMsgs; i++ {
		res := pub.Publish(ctx, &pubsub.Message{
			Data: []byte(fmt.Sprintf("message#%d", i)),
		})
		results = append(results, res)
	}

	for _, r := range results {
		if _, err := r.Get(ctx); err != nil {
			return fmt.Errorf("get publish result: %w", err)
		}
	}

	return nil
}

func TestGenerateCloudEvent(t *testing.T) {
	publishTime := time.Now().UTC()

	tests := []struct {
		name    string
		msg     *pubsub.Message
		wantErr bool
		check   func(t *testing.T, ev cloudevents.Event)
	}{
		{
			name: "plain pubsub message",
			msg: &pubsub.Message{
				ID:          "msg-1",
				Data:        []byte(`{"key":"value"}`),
				PublishTime: publishTime,
				Attributes:  map[string]string{"content-type": "application/json"},
			},
			check: func(t *testing.T, ev cloudevents.Event) {
				if ev.ID() == "" {
					t.Error("expected non-empty ID")
				}
				if ev.Type() != "pubsub.message" {
					t.Errorf("expected type 'pubsub.message', got %q", ev.Type())
				}
				if ev.Source() != "pubsub://subscription-test" {
					t.Errorf("unexpected source: %q", ev.Source())
				}
				if !ev.Time().Equal(publishTime) {
					t.Errorf("expected time %v, got %v", publishTime, ev.Time())
				}
				if string(ev.Data()) != `{"key":"value"}` {
					t.Errorf("unexpected data: %s", string(ev.Data()))
				}
			},
		},
		{
			name: "cloudevents formatted message",
			msg: &pubsub.Message{
				ID:          "msg-2",
				Data:        []byte(`{"data":"content"}`),
				PublishTime: publishTime,
				Attributes: map[string]string{
					"content-type":    "application/json",
					"ce_specversion":  "1.0",
					"ce_id":           "ce-123",
					"ce_source":       "//pubsub.googleapis.com/projects/my-project/topics/my-topic",
					"ce_type":         "com.example.event",
					"ce_time":         publishTime.Format(time.RFC3339),
					"ce_subject":      "test-subject",
					"custom-attr":     "custom-value",
				},
			},
			check: func(t *testing.T, ev cloudevents.Event) {
				if ev.ID() != "ce-123" {
					t.Errorf("expected id 'ce-123', got %q", ev.ID())
				}
				if ev.Type() != "com.example.event" {
					t.Errorf("expected type 'com.example.event', got %q", ev.Type())
				}
				if ev.Source() != "//pubsub.googleapis.com/projects/my-project/topics/my-topic" {
					t.Errorf("unexpected source: %q", ev.Source())
				}
				if ev.SpecVersion() != "1.0" {
					t.Errorf("expected specversion '1.0', got %q", ev.SpecVersion())
				}
				if ev.Subject() != "test-subject" {
					t.Errorf("expected subject 'test-subject', got %q", ev.Subject())
				}
				customAttr, ok := ev.Extensions()["custom-attr"]
				if !ok || customAttr != "custom-value" {
					t.Error("custom-attr extension not found or wrong value")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev, err := generateCloudEvent(tt.msg, "subscription-test")
			if (err != nil) != tt.wantErr {
				t.Fatalf("generateCloudEvent() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if tt.check != nil {
				tt.check(t, ev)
			}
		})
	}
}

func TestProcessMessageRetry(t *testing.T) {
	t.Run("retry exhausted", func(t *testing.T) {
		client, ctx := setup(t)

		subName, err := setupTopicAndSubscription(client, ctx, 3, "retry-exhausted")
		if err != nil {
			t.Fatalf("Failed to setup test: %v", err)
		}

		var attempts atomic.Int32
		h := Handler{
			sleep: 10 * time.Millisecond,
			shouldSucceed: func() bool {
				attempts.Add(1)
				return false
			},
		}

		sub := NewSubscriber[cloudevents.Event](client, h.handle, subName, &Options{
			Concurrency: 1,
			RetryLimit:  2,
			Backoff:     false,
		})

		ctxTimeout, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		_ = sub.Subscribe(ctxTimeout)

		if got := attempts.Load(); got < 3 {
			t.Errorf("expected at least 3 attempts (initial + 2 retries), got %d", got)
		}
	})

	t.Run("retry succeeds", func(t *testing.T) {
		client, ctx := setup(t)

		subName, err := setupTopicAndSubscription(client, ctx, 3, "retry-succeeds")
		if err != nil {
			t.Fatalf("Failed to setup test: %v", err)
		}

		var calls atomic.Int32
		h := Handler{
			sleep: 10 * time.Millisecond,
			shouldSucceed: func() bool {
				return calls.Add(1) > 2
			},
		}

		sub := NewSubscriber[cloudevents.Event](client, h.handle, subName, &Options{
			Concurrency: 1,
			RetryLimit:  3,
			Backoff:     false,
		})

		ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		_ = sub.Subscribe(ctxTimeout)

		if got := h.received; got < 1 {
			t.Errorf("handler should have succeeded at least once, got %d", got)
		}
	})
}
