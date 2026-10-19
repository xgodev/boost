package nats

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	v2 "github.com/cloudevents/sdk-go/v2"
	natsserver "github.com/nats-io/nats-server/v2/test"
	natsgo "github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/wrapper/log"
)

// entry is one line a recorder received: its level and everything printed on it.
type entry struct {
	level string
	text  string
}

// recorder is a log.Logger that keeps every line, with its fields, in memory.
type recorder struct {
	mu      *sync.Mutex
	entries *[]entry
	fields  log.Fields
}

func newRecorder() *recorder {
	return &recorder{mu: &sync.Mutex{}, entries: &[]entry{}, fields: log.Fields{}}
}

func (r *recorder) add(level, msg string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	*r.entries = append(*r.entries, entry{level: level, text: fmt.Sprintf("%s %v", msg, r.fields)})
}

func (r *recorder) lines() []entry {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]entry(nil), *r.entries...)
}

func (r *recorder) with(f log.Fields) log.Logger {
	merged := log.Fields{}
	for k, v := range r.fields {
		merged[k] = v
	}
	for k, v := range f {
		merged[k] = v
	}
	return &recorder{mu: r.mu, entries: r.entries, fields: merged}
}

func (r *recorder) Printf(format string, args ...interface{}) {
	r.add("print", fmt.Sprintf(format, args...))
}
func (r *recorder) Tracef(format string, args ...interface{}) {
	r.add("trace", fmt.Sprintf(format, args...))
}
func (r *recorder) Trace(args ...interface{}) { r.add("trace", fmt.Sprint(args...)) }
func (r *recorder) Debugf(format string, args ...interface{}) {
	r.add("debug", fmt.Sprintf(format, args...))
}
func (r *recorder) Debug(args ...interface{}) { r.add("debug", fmt.Sprint(args...)) }
func (r *recorder) Infof(format string, args ...interface{}) {
	r.add("info", fmt.Sprintf(format, args...))
}
func (r *recorder) Info(args ...interface{}) { r.add("info", fmt.Sprint(args...)) }
func (r *recorder) Warnf(format string, args ...interface{}) {
	r.add("warn", fmt.Sprintf(format, args...))
}
func (r *recorder) Warn(args ...interface{}) { r.add("warn", fmt.Sprint(args...)) }
func (r *recorder) Errorf(format string, args ...interface{}) {
	r.add("error", fmt.Sprintf(format, args...))
}
func (r *recorder) Error(args ...interface{}) { r.add("error", fmt.Sprint(args...)) }
func (r *recorder) Fatalf(format string, args ...interface{}) {
	r.add("fatal", fmt.Sprintf(format, args...))
}
func (r *recorder) Fatal(args ...interface{}) { r.add("fatal", fmt.Sprint(args...)) }
func (r *recorder) Panicf(format string, args ...interface{}) {
	r.add("panic", fmt.Sprintf(format, args...))
}
func (r *recorder) Panic(args ...interface{}) { r.add("panic", fmt.Sprint(args...)) }

func (r *recorder) WithFields(keyValues map[string]interface{}) log.Logger { return r.with(keyValues) }
func (r *recorder) WithField(key string, value interface{}) log.Logger {
	return r.with(log.Fields{key: value})
}
func (r *recorder) WithError(err error) log.Logger { return r.with(log.Fields{"error": err}) }
func (r *recorder) WithTypeOf(obj interface{}) log.Logger {
	return r.with(log.Fields{"type": fmt.Sprintf("%T", obj)})
}
func (r *recorder) ToContext(ctx context.Context) context.Context { return ctx }
func (r *recorder) FromContext(context.Context) log.Logger        { return r }
func (r *recorder) Output() io.Writer                             { return io.Discard }
func (r *recorder) Fields() log.Fields                            { return r.fields }

func TestPublish_NeverLogsThePayloadAndLogsTheIDAtDebug(t *testing.T) {
	boost.Start()
	opts := natsserver.DefaultTestOptions
	opts.Port = -1
	srv := natsserver.RunServer(&opts)
	defer srv.Shutdown()

	conn, err := natsgo.Connect(srv.ClientURL())
	require.NoError(t, err)
	defer conn.Close()

	received := make(chan *natsgo.Msg, 1)
	sub, err := conn.ChanSubscribe("prices", received)
	require.NoError(t, err)
	defer func() { _ = sub.Unsubscribe() }() // intentional: the server shuts down right after
	require.NoError(t, conn.Flush())

	const secret = "card-4111-1111"
	e := v2.NewEvent()
	e.SetID("evt-209")
	e.SetSource("axis-price")
	e.SetType("price")
	e.SetSubject("prices")
	require.NoError(t, e.SetData(v2.ApplicationJSON, map[string]string{"secret": secret}))

	rec := newRecorder()
	previous := log.GetLogger()
	log.Set(rec)
	defer log.Set(previous)

	out, err := New(conn).Publish(context.Background(), []*v2.Event{&e})
	require.NoError(t, err)
	require.Len(t, out, 1)
	assert.NoError(t, out[0].Error)

	select {
	case msg := <-received:
		assert.Contains(t, string(msg.Data), "evt-209")
	case <-time.After(5 * time.Second):
		t.Fatal("the event was not published")
	}

	encoded := base64.StdEncoding.EncodeToString(e.Data())
	var debugMentionsID bool
	for _, l := range rec.lines() {
		assert.NotContains(t, l.text, secret, "payload logged at %s", l.level)
		assert.NotContains(t, l.text, encoded, "payload logged at %s", l.level)
		assert.Contains(t, []string{"trace", "debug"}, l.level, "a publish logs nothing at info or above, got %q", l.text)
		if strings.Contains(l.text, "evt-209") && strings.Contains(l.text, "prices") {
			debugMentionsID = true
		}
	}
	assert.True(t, debugMentionsID, "no debug line with the event id and subject: %v", rec.lines())
}
