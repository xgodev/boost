package pubsub

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
)

// blocking is a handler that signals when it takes a message, holds it
// until release is closed, and records whether its context was cancelled.
type blocking struct {
	started   chan struct{}
	release   chan struct{}
	handled   atomic.Int32
	cancelled atomic.Bool
	fail      bool
}

func newBlocking() *blocking {
	return &blocking{started: make(chan struct{}, 10), release: make(chan struct{})}
}

func (b *blocking) handle(ctx context.Context, _ cloudevents.Event) (cloudevents.Event, error) {
	b.started <- struct{}{}
	<-b.release
	b.cancelled.Store(ctx.Err() != nil)
	b.handled.Add(1)
	if b.fail {
		return cloudevents.Event{}, errors.New("handler failed")
	}
	return cloudevents.Event{}, nil
}

// subscription creates topic and subscription name on the pstest server of
// client.
func subscription(t *testing.T, ctx context.Context, client *pubsub.Client, name string) {
	t.Helper()
	topic := "projects/local-client/topics/" + name
	_, err := client.TopicAdminClient.CreateTopic(ctx, &pubsubpb.Topic{Name: topic})
	require.NoError(t, err)
	_, err = client.SubscriptionAdminClient.CreateSubscription(ctx, &pubsubpb.Subscription{
		Name: "projects/local-client/subscriptions/" + name, Topic: topic, AckDeadlineSeconds: 10,
	})
	require.NoError(t, err)
}

func publish(t *testing.T, ctx context.Context, client *pubsub.Client, topic, data string) {
	t.Helper()
	pub := client.Publisher(topic)
	defer pub.Stop()
	_, err := pub.Publish(ctx, &pubsub.Message{Data: []byte(data)}).Get(ctx)
	require.NoError(t, err)
}

func options(subscriptions ...string) *Options {
	return &Options{Subscriptions: subscriptions, Concurrency: 1, RetryLimit: 100, BackoffBase: time.Millisecond, MaxBackoff: time.Millisecond, ProcessTimeout: 10 * time.Second}
}

// TestHelperRun_StopsOnTheRunContext cancels the run context while the
// handler holds a message (#205): Run blocks until then, the handler
// finishes with a context the stop did not cancel, and Run returns.
func TestHelperRun_StopsOnTheRunContext(t *testing.T) {
	boost.Start()
	client, ctx := setup(t)
	subscription(t, ctx, client, "orders")
	h := newBlocking()
	helper := NewHelperWithOptions[cloudevents.Event](client, h.handle, options("orders"))

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	ran := make(chan error, 1)
	go func() { ran <- helper.Run(runCtx) }()

	publish(t, ctx, client, "orders", `{"item":"book"}`)
	select {
	case <-h.started:
	case err := <-ran:
		t.Fatalf("Run returned (%v) before the stop", err)
	case <-time.After(5 * time.Second):
		t.Fatal("the handler never took the message")
	}

	cancel()
	select {
	case err := <-ran:
		t.Fatalf("Run returned (%v) while the handler held a message", err)
	case <-time.After(200 * time.Millisecond):
	}

	close(h.release)
	select {
	case err := <-ran:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("Run did not return after the handler finished")
	}
	require.EqualValues(t, 1, h.handled.Load())
	require.False(t, h.cancelled.Load(), "the stop cancelled the context of the handler")
}

// TestHelperRun_StopsRetryingOnTheRunContext: a handler that fails while the
// app stops is not retried; the message goes back to Pub/Sub.
func TestHelperRun_StopsRetryingOnTheRunContext(t *testing.T) {
	boost.Start()
	client, ctx := setup(t)
	subscription(t, ctx, client, "orders")
	h := newBlocking()
	h.fail = true
	helper := NewHelperWithOptions[cloudevents.Event](client, h.handle, options("orders"))

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	ran := make(chan error, 1)
	go func() { ran <- helper.Run(runCtx) }()

	publish(t, ctx, client, "orders", `{"item":"book"}`)
	select {
	case <-h.started:
	case <-time.After(5 * time.Second):
		t.Fatal("the handler never took the message")
	}
	cancel()
	close(h.release)
	select {
	case err := <-ran:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("Run kept retrying after the stop")
	}
	require.EqualValues(t, 1, h.handled.Load())
}

// TestHelperRun_FailsWhenASubscriptionIsMissing: one of the subscriptions
// does not exist (#205). Run stops the others and returns an error naming
// it, so the app fails its boot instead of running without it.
func TestHelperRun_FailsWhenASubscriptionIsMissing(t *testing.T) {
	boost.Start()
	client, ctx := setup(t)
	subscription(t, ctx, client, "orders")
	helper := NewHelperWithOptions[cloudevents.Event](client, newBlocking().handle, options("orders", "missing"))

	ran := make(chan error, 1)
	go func() { ran <- helper.Run(ctx) }()
	select {
	case err := <-ran:
		require.ErrorContains(t, err, "pubsub: could not receive on missing")
	case <-time.After(5 * time.Second):
		t.Fatal("Run kept running without subscription missing")
	}
}

// TestNew_RunsUntilTheCommandContextIsDone: the gcp_pubsub command runs the
// helper, with the subscriptions of config, until the context cobra
// executes it with is done (#205), instead of returning at once.
func TestNew_RunsUntilTheCommandContextIsDone(t *testing.T) {
	boost.Start()
	client, ctx := setup(t)
	subscription(t, ctx, client, "changeme")
	cmd := New[cloudevents.Event](client)(newBlocking().handle)
	require.Equal(t, "gcp_pubsub", cmd.Name())

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	cmd.SetArgs([]string{})
	done := make(chan error, 1)
	go func() { done <- cmd.ExecuteContext(runCtx) }()
	select {
	case err := <-done:
		t.Fatalf("the command returned (%v) before its context was done", err)
	case <-time.After(300 * time.Millisecond):
	}

	cancel()
	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("the command did not stop with its context")
	}
}
