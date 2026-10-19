package nats

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/nats-io/nats-server/v2/server"
	natsserver "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
)

// blocking is a handler that signals when it takes an event and holds it
// until release is closed, and records whether its context was cancelled.
type blocking struct {
	started   chan struct{}
	release   chan struct{}
	handled   atomic.Int32
	cancelled atomic.Bool
}

func newBlocking() *blocking {
	return &blocking{started: make(chan struct{}, 10), release: make(chan struct{})}
}

func (b *blocking) handle(ctx context.Context, _ cloudevents.Event) (*cloudevents.Event, error) {
	b.started <- struct{}{}
	<-b.release
	b.cancelled.Store(ctx.Err() != nil)
	b.handled.Add(1)
	return nil, nil
}

func serve(t *testing.T) (*server.Server, *nats.Conn) {
	t.Helper()
	boost.Start()
	opts := natsserver.DefaultTestOptions
	opts.Port = -1
	srv := natsserver.RunServer(&opts)
	t.Cleanup(srv.Shutdown)
	conn, err := nats.Connect(srv.ClientURL())
	require.NoError(t, err)
	t.Cleanup(conn.Close)
	return srv, conn
}

func event(id string) []byte {
	return []byte(`{"specversion":"1.0","id":"` + id + `","source":"test","type":"order.created"}`)
}

// TestHelperRun_StopsOnTheRunContext cancels the run context while the
// handler holds an event (#204): the handler finishes with a context the
// stop did not cancel, Run returns, and no new event is taken.
func TestHelperRun_StopsOnTheRunContext(t *testing.T) {
	srv, conn := serve(t)
	before := srv.NumSubscriptions()
	h := newBlocking()
	helper := NewHelperWithOptions[*cloudevents.Event](conn, h.handle, &Options{Subjects: []string{"orders"}, Queue: "workers"})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ran := make(chan error, 1)
	go func() { ran <- helper.Run(ctx) }()

	require.Eventually(t, func() bool { return srv.NumSubscriptions() > before }, 5*time.Second, 10*time.Millisecond)
	require.NoError(t, conn.Publish("orders", event("1")))
	select {
	case <-h.started:
	case <-time.After(5 * time.Second):
		t.Fatal("the handler never took the event")
	}

	cancel()
	select {
	case err := <-ran:
		t.Fatalf("Run returned (%v) while the handler held an event", err)
	case <-time.After(200 * time.Millisecond):
	}

	close(h.release)
	select {
	case err := <-ran:
		require.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not return after the handler finished")
	}
	require.EqualValues(t, 1, h.handled.Load())
	require.False(t, h.cancelled.Load(), "the stop cancelled the context of the handler")

	require.NoError(t, conn.Publish("orders", event("2")))
	require.NoError(t, conn.Flush())
	time.Sleep(200 * time.Millisecond)
	require.EqualValues(t, 1, h.handled.Load(), "an event published after Run returned was taken")
}

// TestHelperRun_FailsWhenItCannotSubscribe returns the error of a subject it
// cannot receive on, and stops the others, instead of exiting the process or
// idling (#204).
func TestHelperRun_FailsWhenItCannotSubscribe(t *testing.T) {
	srv, conn := serve(t)
	before := srv.NumSubscriptions()
	helper := NewHelperWithOptions[*cloudevents.Event](conn, newBlocking().handle, &Options{Subjects: []string{"orders", ""}, Queue: "workers"})

	ran := make(chan error, 1)
	go func() { ran <- helper.Run(context.Background()) }()
	select {
	case err := <-ran:
		require.ErrorIs(t, err, nats.ErrBadSubject)
	case <-time.After(2 * time.Second):
		t.Fatal("Run idled on a subject it could not receive on")
	}
	require.NoError(t, conn.Flush())
	require.Eventually(t, func() bool { return srv.NumSubscriptions() == before }, 2*time.Second, 10*time.Millisecond)
}

// TestHelperStart_ReturnsWhenItCannotSubscribe keeps Start for callers
// without a context: it logs the error and returns.
func TestHelperStart_ReturnsWhenItCannotSubscribe(t *testing.T) {
	_, conn := serve(t)
	conn.Close()
	helper := NewHelperWithOptions[*cloudevents.Event](conn, newBlocking().handle, &Options{Subjects: []string{"orders"}, Queue: "workers"})

	done := make(chan struct{})
	go func() { helper.Start(); close(done) }()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Start idled on a subject it could not receive on")
	}
}

// TestNew_RunsUntilTheCommandContextIsDone: the cloudevents_nats command
// runs the helper with the context cobra executes it with.
func TestNew_RunsUntilTheCommandContextIsDone(t *testing.T) {
	_, conn := serve(t)
	cmd := New[*cloudevents.Event](conn)(newBlocking().handle)
	require.Equal(t, "cloudevents_nats", cmd.Name())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	cmd.SetArgs([]string{})
	require.NoError(t, cmd.ExecuteContext(ctx))
}
