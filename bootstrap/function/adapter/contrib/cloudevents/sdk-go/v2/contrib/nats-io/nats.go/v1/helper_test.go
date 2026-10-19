package nats

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	natsserver "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
)

// TestHelperSubscribeCallsAnEventHandler feeds a function that replies with
// an event, as a @Function handler does: the nats consumer has no responder,
// so the reply stays with the function and its middlewares.
func TestHelperSubscribeCallsAnEventHandler(t *testing.T) {
	boost.Start()

	opts := natsserver.DefaultTestOptions
	opts.Port = -1
	srv := natsserver.RunServer(&opts)
	defer srv.Shutdown()

	conn, err := nats.Connect(srv.ClientURL())
	require.NoError(t, err)
	defer conn.Close()

	var calls atomic.Int32
	handler := func(ctx context.Context, in cloudevents.Event) (*cloudevents.Event, error) {
		calls.Add(1)
		return &in, nil
	}
	h := NewHelperWithOptions[*cloudevents.Event](conn, handler, &Options{Subjects: []string{"messages"}, Queue: "test"})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go h.subscribe(ctx, "messages")

	// nats drops what is published before the helper subscribes
	event := []byte(`{"specversion":"1.0","id":"1","source":"test","type":"message.sent"}`)
	require.Eventually(t, func() bool {
		require.NoError(t, conn.Publish("messages", event))
		return calls.Load() > 0
	}, 5*time.Second, 50*time.Millisecond, "the handler was never called")
}
