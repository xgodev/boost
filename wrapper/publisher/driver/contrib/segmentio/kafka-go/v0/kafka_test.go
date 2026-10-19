package kafka

import (
	"context"
	"errors"
	"testing"
	"time"

	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeWriter struct {
	calls [][]kafka.Message
	err   error
}

func (f *fakeWriter) WriteMessages(_ context.Context, msgs ...kafka.Message) error {
	f.calls = append(f.calls, msgs)
	return f.err
}

func event(t *testing.T, id, subject string) *v2.Event {
	t.Helper()
	e := v2.NewEvent()
	e.SetID(id)
	e.SetSource("/orders")
	e.SetType("order.created")
	e.SetSubject(subject)
	e.SetTime(time.Date(2026, 10, 18, 12, 30, 0, 0, time.UTC))
	require.NoError(t, e.SetData(v2.ApplicationJSON, map[string]any{"id": 1}))
	return &e
}

func headerMap(m kafka.Message) map[string]string {
	h := map[string]string{}
	for _, x := range m.Headers {
		h[x.Key] = string(x.Value)
	}
	return h
}

func TestPublish_WritesEachEventInBinaryMode(t *testing.T) {
	w := &fakeWriter{}
	e := event(t, "e1", "orders")
	e.SetExtension("tenant", "acme")

	res, err := newDriver(w).Publish(context.Background(), []*v2.Event{e})
	require.NoError(t, err)
	require.Len(t, res, 1)
	assert.NoError(t, res[0].Error)
	assert.Same(t, e, res[0].Event)

	require.Len(t, w.calls, 1)
	require.Len(t, w.calls[0], 1)
	m := w.calls[0][0]
	assert.Equal(t, "orders", m.Topic)
	assert.Equal(t, "e1", string(m.Key), "the event ID is the default key")
	assert.JSONEq(t, `{"id":1}`, string(m.Value))
	assert.Equal(t, map[string]string{
		"content-type":   "application/json",
		"ce_specversion": "1.0",
		"ce_id":          "e1",
		"ce_source":      "/orders",
		"ce_type":        "order.created",
		"ce_subject":     "orders",
		"ce_time":        "2026-10-18T12:30:00Z",
		"ce_tenant":      "acme",
	}, headerMap(m))
}

func TestPublish_KeyFromExtensions(t *testing.T) {
	w := &fakeWriter{}
	withKey := event(t, "e1", "orders")
	withKey.SetExtension("key", "k1")
	withKey.SetExtension("partitionkey", "ignored")
	withPartitionKey := event(t, "e2", "orders")
	withPartitionKey.SetExtension("partitionkey", "p2")

	_, err := newDriver(w).Publish(context.Background(), []*v2.Event{withKey, withPartitionKey})
	require.NoError(t, err)

	require.Len(t, w.calls[0], 2)
	assert.Equal(t, "k1", string(w.calls[0][0].Key))
	assert.Equal(t, "p2", string(w.calls[0][1].Key))
}

func TestPublish_EventWithoutIDGetsOne(t *testing.T) {
	w := &fakeWriter{}
	e := event(t, "", "orders")

	_, err := newDriver(w).Publish(context.Background(), []*v2.Event{e})
	require.NoError(t, err)
	assert.NotEmpty(t, e.ID())
	assert.Equal(t, e.ID(), string(w.calls[0][0].Key))
}

func TestPublish_EventWithoutSubjectFailsAlone(t *testing.T) {
	w := &fakeWriter{}
	bad := event(t, "e1", "")
	good := event(t, "e2", "orders")

	res, err := newDriver(w).Publish(context.Background(), []*v2.Event{bad, good})
	require.NoError(t, err)
	assert.ErrorContains(t, res[0].Error, "subject")
	assert.NoError(t, res[1].Error)
	require.Len(t, w.calls[0], 1)
	assert.Equal(t, "e2", string(w.calls[0][0].Key))
}

func TestPublish_WriteErrorsMapToTheirEvents(t *testing.T) {
	w := &fakeWriter{err: kafka.WriteErrors{nil, errors.New("leader not available")}}
	bad := event(t, "e0", "")
	ok := event(t, "e1", "orders")
	failed := event(t, "e2", "orders")

	res, err := newDriver(w).Publish(context.Background(), []*v2.Event{bad, ok, failed})
	require.NoError(t, err)
	assert.ErrorContains(t, res[0].Error, "subject")
	assert.NoError(t, res[1].Error)
	assert.ErrorContains(t, res[2].Error, "leader not available")
}

func TestPublish_WriteErrorFailsEverySentEvent(t *testing.T) {
	w := &fakeWriter{err: errors.New("broker down")}

	res, err := newDriver(w).Publish(context.Background(), []*v2.Event{event(t, "e1", "a"), event(t, "e2", "b")})
	require.NoError(t, err)
	assert.ErrorContains(t, res[0].Error, "broker down")
	assert.ErrorContains(t, res[1].Error, "broker down")
}

func TestPublish_NothingToWrite(t *testing.T) {
	w := &fakeWriter{}

	res, err := newDriver(w).Publish(context.Background(), []*v2.Event{event(t, "e1", "")})
	require.NoError(t, err)
	assert.Error(t, res[0].Error)
	assert.Empty(t, w.calls)
}
