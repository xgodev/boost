package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/bootstrap/function/examples/nats/gen/builder/worker"
	"github.com/xgodev/boost/bootstrap/function/examples/nats/internal/app/worker/nats/handler/message/replier"
	"github.com/xgodev/boost/wrapper/publisher"
	"go.uber.org/fx"
)

var url string

// TestMain starts an embedded nats server, points the nats factory at it and
// loads configs/config.yaml, which makes the adapter the default command.
func TestMain(m *testing.M) {
	ns, err := server.NewServer(&server.Options{Host: "127.0.0.1", Port: -1})
	if err != nil {
		panic(err)
	}
	go ns.Start()
	if !ns.ReadyForConnections(10 * time.Second) {
		panic("nats server not ready")
	}
	url = ns.ClientURL()
	os.Setenv("BOOST_FACTORY_NATS_URL", url)
	os.Setenv("CONF", "../../configs/config.yaml")
	boost.Start()
	code := m.Run()
	ns.Shutdown()
	os.Exit(code)
}

// recorder keeps the events the publisher middleware publishes.
type recorder struct {
	publisher.Driver
	mu     sync.Mutex
	events []*cloudevents.Event
}

func (r *recorder) Publish(ctx context.Context, events []*cloudevents.Event) ([]publisher.PublishOutput, error) {
	r.mu.Lock()
	r.events = append(r.events, events...)
	r.mu.Unlock()
	return r.Driver.Publish(ctx, events)
}

func (r *recorder) got() []*cloudevents.Event {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]*cloudevents.Event(nil), r.events...)
}

func event(message string) []byte {
	return []byte(fmt.Sprintf(`{"specversion":"1.0","id":"1","source":"test","type":"message.sent","datacontenttype":"application/json","data":{"message":%q}}`, message))
}

// TestWorker starts app worker and publishes events on subject messages.
func TestWorker(t *testing.T) {
	// the function runs the default command, nats, when the command line
	// names none: drop the flags of go test.
	os.Args = os.Args[:1]

	rec := &recorder{}
	app := fx.New(append(worker.Modules(), fx.NopLogger,
		fx.Decorate(func(d publisher.Driver) publisher.Driver { rec.Driver = d; return rec }),
	)...)
	require.NoError(t, app.Err())
	require.NoError(t, app.Start(context.Background()))

	conn, err := nats.Connect(url)
	require.NoError(t, err)
	defer conn.Close()

	// nats drops what is published before the adapter subscribes; an empty
	// message is refused, so it has no reply
	require.Eventually(t, func() bool {
		require.NoError(t, conn.Publish("messages", event("")))
		require.NoError(t, conn.Publish("messages", event("hi")))
		return len(rec.got()) > 0
	}, 10*time.Second, 100*time.Millisecond, "the publisher middleware published no reply")

	for _, e := range rec.got() {
		assert.Equal(t, replier.RepliesSubject, e.Subject())
		assert.JSONEq(t, `{"message":"hello, hi"}`, string(e.Data()))
	}
}
