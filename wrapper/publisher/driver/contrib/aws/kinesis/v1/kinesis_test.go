package kinesis

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awskinesis "github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"
	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeKinesis records single and bulk puts per stream and fails the streams
// in failing.
type fakeKinesis struct {
	mu      sync.Mutex
	singles []*awskinesis.PutRecordInput
	bulks   map[string][]types.PutRecordsRequestEntry
	failing map[string]bool
}

func (f *fakeKinesis) Publish(_ context.Context, input *awskinesis.PutRecordInput) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.singles = append(f.singles, input)
	if f.failing[aws.ToString(input.StreamName)] {
		return fmt.Errorf("stream %s is down", aws.ToString(input.StreamName))
	}
	return nil
}

func (f *fakeKinesis) BulkPublish(_ context.Context, entries []types.PutRecordsRequestEntry, stream string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.bulks == nil {
		f.bulks = map[string][]types.PutRecordsRequestEntry{}
	}
	f.bulks[stream] = append(f.bulks[stream], entries...)
	if f.failing[stream] {
		return fmt.Errorf("stream %s is down", stream)
	}
	return nil
}

func event(id, subject string, exts map[string]string) *v2.Event {
	e := v2.NewEvent()
	e.SetID(id)
	e.SetSource("/orders")
	e.SetType("order.created")
	e.SetSubject(subject)
	for k, v := range exts {
		e.SetExtension(k, v)
	}
	_ = e.SetData(v2.ApplicationJSON, map[string]string{"id": id})
	return &e
}

func TestPublish_OneEventOfAStreamIsAPutRecord(t *testing.T) {
	f := &fakeKinesis{}
	e := event("1", "orders", map[string]string{"partitionkey": "customer-7", "target": "data"})

	out, err := New(f).Publish(context.Background(), []*v2.Event{e})
	require.NoError(t, err)
	require.Len(t, out, 1)
	assert.Same(t, e, out[0].Event)
	assert.NoError(t, out[0].Error)

	require.Len(t, f.singles, 1)
	assert.Equal(t, "orders", aws.ToString(f.singles[0].StreamName))
	assert.Equal(t, "customer-7", aws.ToString(f.singles[0].PartitionKey))
	assert.JSONEq(t, `{"id":"1"}`, string(f.singles[0].Data))
	assert.Empty(t, f.bulks)
}

func TestPublish_EventsOfAStreamGoInOnePutRecords(t *testing.T) {
	f := &fakeKinesis{}
	events := []*v2.Event{
		event("1", "orders", nil),
		event("2", "payments", nil),
		event("3", "orders", map[string]string{"partitionkey": "customer-7"}),
	}

	out, err := New(f).Publish(context.Background(), events)
	require.NoError(t, err)
	for i, o := range out {
		assert.Same(t, events[i], o.Event)
		assert.NoError(t, o.Error)
	}

	require.Len(t, f.bulks["orders"], 2)
	assert.Equal(t, "1", aws.ToString(f.bulks["orders"][0].PartitionKey), "without the extension the key is the event ID")
	assert.Equal(t, "customer-7", aws.ToString(f.bulks["orders"][1].PartitionKey))
	require.Len(t, f.singles, 1)
	assert.Equal(t, "payments", aws.ToString(f.singles[0].StreamName))
}

func TestPublish_AFailingStreamFailsOnlyItsEvents(t *testing.T) {
	f := &fakeKinesis{failing: map[string]bool{"orders": true}}
	events := []*v2.Event{event("1", "orders", nil), event("2", "payments", nil), event("3", "orders", nil)}

	out, err := New(f).Publish(context.Background(), events)
	require.NoError(t, err)
	assert.ErrorContains(t, out[0].Error, "stream orders is down")
	assert.NoError(t, out[1].Error)
	assert.ErrorContains(t, out[2].Error, "stream orders is down")
}

func TestPublish_AnEventWithoutIDGetsAPartitionKey(t *testing.T) {
	f := &fakeKinesis{}
	e := event("", "orders", nil)

	out, err := New(f).Publish(context.Background(), []*v2.Event{e})
	require.NoError(t, err)
	require.NoError(t, out[0].Error)
	assert.NotEmpty(t, aws.ToString(f.singles[0].PartitionKey))
}
