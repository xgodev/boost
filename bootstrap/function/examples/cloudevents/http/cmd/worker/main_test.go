package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/bootstrap/function/examples/cloudevents/http/gen/builder/worker"
	"github.com/xgodev/boost/bootstrap/function/examples/cloudevents/http/internal/app/worker/cloudevents/handler/message/replier"
	"github.com/xgodev/boost/wrapper/publisher"
	"go.uber.org/fx"
)

var url string

// TestMain puts the adapter on a free port and loads configs/config.yaml,
// which makes it the default command.
func TestMain(m *testing.M) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()
	url = fmt.Sprintf("http://127.0.0.1:%d/", port)
	os.Setenv("BOOST_BOOTSTRAP_FUNCTION_ADAPTER_CLOUDEVENTS_HTTP_PORT", strconv.Itoa(port))
	os.Setenv("CONF", "../../configs/config.yaml")
	boost.Start()
	os.Exit(m.Run())
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

func post(message string) (*http.Response, error) {
	body := fmt.Sprintf(`{"specversion":"1.0","id":"1","source":"test","type":"message.sent","datacontenttype":"application/json","data":{"message":%q}}`, message)
	return http.Post(url, "application/cloudevents+json", strings.NewReader(body))
}

// TestWorker starts app worker and posts events to the adapter.
func TestWorker(t *testing.T) {
	// the function runs the default command, cloudevents_http, when the
	// command line names none: drop the flags of go test.
	os.Args = os.Args[:1]

	rec := &recorder{}
	app := fx.New(append(worker.Modules(), fx.NopLogger,
		fx.Decorate(func(d publisher.Driver) publisher.Driver { rec.Driver = d; return rec }),
	)...)
	require.NoError(t, app.Err())
	require.NoError(t, app.Start(context.Background()))

	var res *http.Response
	require.Eventually(t, func() bool {
		var err error
		res, err = post("hi")
		return err == nil
	}, 10*time.Second, 50*time.Millisecond, "the adapter did not listen")
	res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	events := rec.got()
	require.Len(t, events, 1, "the publisher middleware publishes the reply")
	assert.Equal(t, replier.RepliesSubject, events[0].Subject())
	assert.JSONEq(t, `{"message":"hello, hi"}`, string(events[0].Data()))

	res, err := post("")
	require.NoError(t, err)
	res.Body.Close()
	assert.Equal(t, http.StatusBadRequest, res.StatusCode, "an empty message is refused")
	assert.Len(t, rec.got(), 1, "a refused message has no reply")
}
