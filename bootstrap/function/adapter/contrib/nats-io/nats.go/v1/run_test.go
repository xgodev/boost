package nats

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/nats-io/nats-server/v2/server"
	natsserver "github.com/nats-io/nats-server/v2/test"
	n "github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
)

// blocking is a handler that signals when it takes a message and holds it
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

func (b *blocking) handle(ctx context.Context, _ v2.Event) (*v2.Event, error) {
	b.started <- struct{}{}
	<-b.release
	b.cancelled.Store(ctx.Err() != nil)
	b.handled.Add(1)
	return nil, nil
}

func serve(t *testing.T) (*server.Server, *n.Conn) {
	t.Helper()
	boost.Start()
	opts := natsserver.DefaultTestOptions
	opts.Port = -1
	srv := natsserver.RunServer(&opts)
	t.Cleanup(srv.Shutdown)
	conn, err := n.Connect(srv.ClientURL())
	require.NoError(t, err)
	t.Cleanup(conn.Close)
	return srv, conn
}

// TestHelperRun_StopsOnTheRunContext cancels the run context while the
// handler holds a message (#204): the handler finishes with a context the
// stop did not cancel, Run returns, and no new message is taken.
func TestHelperRun_StopsOnTheRunContext(t *testing.T) {
	srv, conn := serve(t)
	before := srv.NumSubscriptions()
	h := newBlocking()
	helper := NewHelperWithOptions[*v2.Event](conn, h.handle, &Options{Subjects: []string{"orders"}, Queue: "workers"})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ran := make(chan error, 1)
	go func() { ran <- helper.Run(ctx) }()

	require.Eventually(t, func() bool { return srv.NumSubscriptions() > before }, 5*time.Second, 10*time.Millisecond)
	require.NoError(t, conn.Publish("orders", []byte(`{"item":"book"}`)))
	select {
	case <-h.started:
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
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not return after the handler finished")
	}
	require.EqualValues(t, 1, h.handled.Load())
	require.False(t, h.cancelled.Load(), "the stop cancelled the context of the handler")

	require.NoError(t, conn.Publish("orders", []byte(`{"item":"pen"}`)))
	require.NoError(t, conn.Flush())
	time.Sleep(200 * time.Millisecond)
	require.EqualValues(t, 1, h.handled.Load(), "a message published after Run returned was taken")
}

// TestHelperRun_FailsWhenItCannotSubscribe returns the subscribe error
// instead of logging it and idling (#204).
func TestHelperRun_FailsWhenItCannotSubscribe(t *testing.T) {
	_, conn := serve(t)
	conn.Close()
	helper := NewHelperWithOptions[*v2.Event](conn, newBlocking().handle, &Options{Subjects: []string{"orders"}, Queue: "workers"})

	ran := make(chan error, 1)
	go func() { ran <- helper.Run(context.Background()) }()
	select {
	case err := <-ran:
		require.ErrorIs(t, err, n.ErrConnectionClosed)
	case <-time.After(2 * time.Second):
		t.Fatal("Run idled on a subscription that failed")
	}
}

// TestHelperRun_UnsubscribesWhatItSubscribed: a subject that fails after
// another one subscribed leaves no subscription behind.
func TestHelperRun_UnsubscribesWhatItSubscribed(t *testing.T) {
	srv, conn := serve(t)
	before := srv.NumSubscriptions()
	helper := NewHelperWithOptions[*v2.Event](conn, newBlocking().handle, &Options{Subjects: []string{"orders", ""}, Queue: "workers"})

	require.ErrorIs(t, helper.Run(context.Background()), n.ErrBadSubject)
	require.NoError(t, conn.Flush())
	require.Eventually(t, func() bool { return srv.NumSubscriptions() == before }, 2*time.Second, 10*time.Millisecond)
}

// TestHelperStart_ReturnsWhenItCannotSubscribe keeps Start for callers
// without a context: it logs the error and returns.
func TestHelperStart_ReturnsWhenItCannotSubscribe(t *testing.T) {
	_, conn := serve(t)
	conn.Close()
	helper := NewHelperWithOptions[*v2.Event](conn, newBlocking().handle, &Options{Subjects: []string{"orders"}, Queue: "workers"})

	done := make(chan struct{})
	go func() { helper.Start(); close(done) }()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Start idled on a subscription that failed")
	}
}

// TestHelperRun_ReturnsWhenTheConnectionClosed: a connection closed under Run
// leaves nothing to drain, and Run still returns once ctx is done.
func TestHelperRun_ReturnsWhenTheConnectionClosed(t *testing.T) {
	srv, conn := serve(t)
	before := srv.NumSubscriptions()
	helper := NewHelperWithOptions[*v2.Event](conn, newBlocking().handle, &Options{Subjects: []string{"orders"}, Queue: "workers"})

	ctx, cancel := context.WithCancel(context.Background())
	ran := make(chan error, 1)
	go func() { ran <- helper.Run(ctx) }()
	require.Eventually(t, func() bool { return srv.NumSubscriptions() > before }, 5*time.Second, 10*time.Millisecond)

	conn.Close()
	cancel()
	select {
	case err := <-ran:
		require.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not return")
	}
}

// TestNew_RunsUntilTheCommandContextIsDone: the nats command runs the helper
// with the context cobra executes it with.
func TestNew_RunsUntilTheCommandContextIsDone(t *testing.T) {
	_, conn := serve(t)
	cmd := New[*v2.Event](conn)(newBlocking().handle)
	require.Equal(t, "nats", cmd.Name())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	cmd.SetArgs([]string{})
	require.NoError(t, cmd.ExecuteContext(ctx))
}
