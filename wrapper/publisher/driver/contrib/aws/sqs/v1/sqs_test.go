package sqs

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awssqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeSQS resolves a queue name to "https://sqs/<name>", fails the queues in
// unknown and records every message.
type fakeSQS struct {
	mu      sync.Mutex
	inputs  []*awssqs.SendMessageInput
	unknown map[string]bool
}

func (f *fakeSQS) Publish(_ context.Context, input *awssqs.SendMessageInput) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.inputs = append(f.inputs, input)
	return nil
}

func (f *fakeSQS) ResolveQueueUrl(_ context.Context, name string) (*string, error) {
	if f.unknown[name] {
		return nil, fmt.Errorf("queue %s does not exist", name)
	}
	return aws.String("https://sqs/" + name), nil
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

func TestPublish_SendsEachEventToTheQueueOfItsSubject(t *testing.T) {
	f := &fakeSQS{}
	events := []*v2.Event{
		event("1", "orders.fifo", map[string]string{"target": "data", "group": "customer-7"}),
		event("2", "payments", nil),
	}

	out, err := New(f).Publish(context.Background(), events)
	require.NoError(t, err)
	require.Len(t, out, 2)
	for i, o := range out {
		assert.Same(t, events[i], o.Event)
		assert.NoError(t, o.Error)
	}

	byQueue := map[string]*awssqs.SendMessageInput{}
	for _, in := range f.inputs {
		byQueue[aws.ToString(in.QueueUrl)] = in
	}
	require.Len(t, byQueue, 2)

	orders := byQueue["https://sqs/orders.fifo"]
	assert.JSONEq(t, `{"id":"1"}`, aws.ToString(orders.MessageBody))
	assert.Equal(t, "customer-7", aws.ToString(orders.MessageGroupId))

	payments := byQueue["https://sqs/payments"]
	assert.Contains(t, aws.ToString(payments.MessageBody), `"type":"order.created"`)
	assert.Nil(t, payments.MessageGroupId)
}

func TestPublish_AnUnknownQueueFailsOnlyItsEvent(t *testing.T) {
	f := &fakeSQS{unknown: map[string]bool{"missing": true}}
	events := []*v2.Event{event("1", "orders", nil), event("2", "missing", nil)}

	out, err := New(f).Publish(context.Background(), events)
	require.NoError(t, err)
	assert.NoError(t, out[0].Error)
	assert.ErrorContains(t, out[1].Error, "queue missing does not exist")
	assert.Len(t, f.inputs, 1)
}
