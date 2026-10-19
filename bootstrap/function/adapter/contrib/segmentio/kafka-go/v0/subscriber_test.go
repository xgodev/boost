package kafka

import (
	"context"
	"errors"
	"io"
	"sync"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeReader serves msgs, then returns err; with no err it blocks until ctx
// ends, like a reader with nothing to fetch.
type fakeReader struct {
	mu        sync.Mutex
	msgs      []kafka.Message
	err       error
	committed []kafka.Message
	commitErr error
}

func (f *fakeReader) FetchMessage(ctx context.Context) (kafka.Message, error) {
	f.mu.Lock()
	if len(f.msgs) > 0 {
		m := f.msgs[0]
		f.msgs = f.msgs[1:]
		f.mu.Unlock()
		return m, nil
	}
	err := f.err
	f.mu.Unlock()
	if err != nil {
		return kafka.Message{}, err
	}
	<-ctx.Done()
	return kafka.Message{}, ctx.Err()
}

func (f *fakeReader) CommitMessages(_ context.Context, msgs ...kafka.Message) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.committed = append(f.committed, msgs...)
	return f.commitErr
}

func (f *fakeReader) commits() []kafka.Message {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]kafka.Message(nil), f.committed...)
}

// recorder is a handler that keeps what it received and fails while fail
// says so.
type recorder struct {
	mu     sync.Mutex
	events []cloudevents.Event
	fail   func(calls int) bool
}

func (r *recorder) handle(_ context.Context, in cloudevents.Event) (*event.Event, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = append(r.events, in)
	if r.fail != nil && r.fail(len(r.events)) {
		return nil, errors.New("handler failed")
	}
	return nil, nil
}

func (r *recorder) received() []cloudevents.Event {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]cloudevents.Event(nil), r.events...)
}

func options() *Options {
	return &Options{
		Topics:      []string{"orders"},
		Buffer:      10,
		Backoff:     true,
		BackoffBase: time.Millisecond,
		MaxBackoff:  2 * time.Millisecond,
		RetryLimit:  -1,
	}
}

func subscribe(t *testing.T, r *fakeReader, h *recorder, o *Options) {
	t.Helper()
	r.err = io.EOF // a closed reader ends Subscribe once every message is handled
	require.NoError(t, NewSubscriber[*event.Event](r, h.handle, o).Subscribe(context.Background()))
}

func header(k, v string) kafka.Header { return kafka.Header{Key: k, Value: []byte(v)} }

func TestSubscribe_BinaryMode(t *testing.T) {
	msg := kafka.Message{
		Topic: "orders", Partition: 0, Offset: 3,
		Value: []byte(`{"id":1}`),
		Headers: []kafka.Header{
			header("ce_specversion", "1.0"),
			header("ce_id", "e1"),
			header("ce_source", "/orders"),
			header("ce_type", "order.created"),
			header("ce_subject", "orders"),
			header("ce_time", "2026-10-18T12:30:00.5Z"),
			header("ce_dataschema", "https://schema/order"),
			header("ce_tenant", "acme"),
			header("content-type", "application/json"),
			header("traceparent", "ignored: not a ce_ header"),
		},
	}
	r := &fakeReader{msgs: []kafka.Message{msg}}
	h := &recorder{}

	subscribe(t, r, h, options())

	require.Len(t, h.received(), 1)
	in := h.received()[0]
	assert.Equal(t, "1.0", in.SpecVersion())
	assert.Equal(t, "e1", in.ID())
	assert.Equal(t, "/orders", in.Source())
	assert.Equal(t, "order.created", in.Type())
	assert.Equal(t, "orders", in.Subject())
	assert.Equal(t, time.Date(2026, 10, 18, 12, 30, 0, 5e8, time.UTC), in.Time().UTC())
	assert.Equal(t, "https://schema/order", in.DataSchema())
	assert.Equal(t, "application/json", in.DataContentType())
	assert.Equal(t, map[string]any{"tenant": "acme"}, in.Extensions())
	assert.JSONEq(t, `{"id":1}`, string(in.Data()))

	require.Len(t, r.commits(), 1)
	assert.Equal(t, int64(3), r.commits()[0].Offset)
}

func TestSubscribe_StructuredMode(t *testing.T) {
	value := `{"specversion":"1.0","id":"e2","source":"/s","type":"t","datacontenttype":"application/json","data":{"a":1}}`
	withType := kafka.Message{Topic: "orders", Value: []byte(value), Headers: []kafka.Header{header("content-type", "application/cloudevents+json")}}
	withoutType := kafka.Message{Topic: "orders", Offset: 1, Value: []byte(value)}
	r := &fakeReader{msgs: []kafka.Message{withType, withoutType}}
	h := &recorder{}

	subscribe(t, r, h, options())

	require.Len(t, h.received(), 2)
	for _, in := range h.received() {
		assert.Equal(t, "e2", in.ID())
		assert.Equal(t, "/s", in.Source())
		assert.JSONEq(t, `{"a":1}`, string(in.Data()))
	}
	assert.Len(t, r.commits(), 2)
}

func TestSubscribe_RawValue(t *testing.T) {
	at := time.Date(2026, 10, 18, 9, 0, 0, 0, time.UTC)
	msg := kafka.Message{Topic: "orders", Partition: 2, Offset: 7, Time: at, Value: []byte(`{"a":1}`)}
	r := &fakeReader{msgs: []kafka.Message{msg}}
	h := &recorder{}

	subscribe(t, r, h, options())

	require.Len(t, h.received(), 1)
	in := h.received()[0]
	assert.Equal(t, "orders-2-7", in.ID(), "topic, partition and offset identify the record")
	assert.Equal(t, "kafka://orders/2", in.Source())
	assert.Equal(t, "kafka.message", in.Type())
	assert.Equal(t, "orders", in.Subject())
	assert.Equal(t, at, in.Time().UTC())
	assert.Equal(t, "application/json", in.DataContentType())
	assert.JSONEq(t, `{"a":1}`, string(in.Data()))
	assert.Len(t, r.commits(), 1)
}

func TestSubscribe_UndecodableMessageIsCommittedAndSkipped(t *testing.T) {
	bad := kafka.Message{Topic: "orders", Value: []byte("not json"), Headers: []kafka.Header{header("content-type", "application/cloudevents+json")}}
	r := &fakeReader{msgs: []kafka.Message{bad}}
	h := &recorder{}

	subscribe(t, r, h, options())

	assert.Empty(t, h.received())
	assert.Len(t, r.commits(), 1)
}

func TestSubscribe_RetriesUntilTheHandlerSucceeds(t *testing.T) {
	r := &fakeReader{msgs: []kafka.Message{{Topic: "orders", Value: []byte(`{}`)}}}
	h := &recorder{fail: func(calls int) bool { return calls < 3 }}

	subscribe(t, r, h, options())

	assert.Len(t, h.received(), 3)
	assert.Len(t, r.commits(), 1, "committed once, after the success")
}

func TestSubscribe_RetryLimitSkipsTheMessage(t *testing.T) {
	r := &fakeReader{msgs: []kafka.Message{{Topic: "orders", Value: []byte(`{}`)}}}
	h := &recorder{fail: func(int) bool { return true }}
	o := options()
	o.RetryLimit = 2

	subscribe(t, r, h, o)

	assert.Len(t, h.received(), 3, "the first attempt and 2 retries")
	assert.Len(t, r.commits(), 1, "a message given up on is committed, or it would block its partition")
}

func TestSubscribe_KeepsTheOrderOfEachPartition(t *testing.T) {
	var msgs []kafka.Message
	for offset := int64(0); offset < 5; offset++ {
		for _, p := range []int{0, 1} {
			msgs = append(msgs, kafka.Message{Topic: "orders", Partition: p, Offset: offset, Value: []byte(`{}`)})
		}
	}
	r := &fakeReader{msgs: msgs}
	h := &recorder{}

	subscribe(t, r, h, options())

	order := map[string][]string{}
	for _, in := range h.received() {
		order[in.Source()] = append(order[in.Source()], in.ID())
	}
	assert.Equal(t, []string{"orders-0-0", "orders-0-1", "orders-0-2", "orders-0-3", "orders-0-4"}, order["kafka://orders/0"])
	assert.Equal(t, []string{"orders-1-0", "orders-1-1", "orders-1-2", "orders-1-3", "orders-1-4"}, order["kafka://orders/1"])
	assert.Len(t, r.commits(), 10)
}

func TestSubscribe_CommitErrorDoesNotStopTheSubscriber(t *testing.T) {
	r := &fakeReader{
		msgs:      []kafka.Message{{Topic: "orders", Value: []byte(`{}`)}, {Topic: "orders", Offset: 1, Value: []byte(`{}`)}},
		commitErr: errors.New("rebalance in progress"),
	}
	h := &recorder{}

	subscribe(t, r, h, options())

	assert.Len(t, h.received(), 2)
}

func TestSubscribe_FetchErrorIsReturned(t *testing.T) {
	r := &fakeReader{err: errors.New("unknown topic")}

	err := NewSubscriber[*event.Event](r, (&recorder{}).handle, options()).Subscribe(context.Background())
	assert.ErrorContains(t, err, "unknown topic")
}

func TestSubscribe_StopsWhenTheContextEnds(t *testing.T) {
	r := &fakeReader{msgs: []kafka.Message{{Topic: "orders", Value: []byte(`{}`)}}}
	h := &recorder{fail: func(int) bool { return true }}
	o := options()
	o.BackoffBase = time.Hour // a retry waiting on its backoff
	o.MaxBackoff = time.Hour

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- NewSubscriber[*event.Event](r, h.handle, o).Subscribe(ctx) }()

	require.Eventually(t, func() bool { return len(h.received()) == 1 }, time.Second, time.Millisecond)
	cancel()

	select {
	case err := <-done:
		assert.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("Subscribe did not return after the context ended")
	}
	assert.Empty(t, r.commits(), "an unfinished message stays uncommitted, to be redelivered")
}
