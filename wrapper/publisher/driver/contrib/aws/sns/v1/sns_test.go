package sns

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awssns "github.com/aws/aws-sdk-go-v2/service/sns"
	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeSNS records every publish and fails the topics in failing.
type fakeSNS struct {
	mu      sync.Mutex
	inputs  []*awssns.PublishInput
	failing map[string]bool
}

func (f *fakeSNS) Publish(_ context.Context, input *awssns.PublishInput) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.inputs = append(f.inputs, input)
	if f.failing[aws.ToString(input.TopicArn)] {
		return fmt.Errorf("topic %s is down", aws.ToString(input.TopicArn))
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

func TestPublish_SendsEachEventToTheTopicOfItsSubject(t *testing.T) {
	f := &fakeSNS{}
	events := []*v2.Event{
		event("1", "arn:aws:sns:sa-east-1:1:orders", map[string]string{"target": "data"}),
		event("2", "arn:aws:sns:sa-east-1:1:payments", nil),
	}

	out, err := New(f).Publish(context.Background(), events)
	require.NoError(t, err)
	require.Len(t, out, 2)
	for i, o := range out {
		assert.Same(t, events[i], o.Event)
		assert.NoError(t, o.Error)
	}

	require.Len(t, f.inputs, 2)
	byTopic := map[string]*awssns.PublishInput{}
	for _, in := range f.inputs {
		assert.Equal(t, "json", aws.ToString(in.MessageStructure))
		byTopic[aws.ToString(in.TopicArn)] = in
	}

	var msg struct{ Default string }
	require.NoError(t, json.Unmarshal([]byte(aws.ToString(byTopic["arn:aws:sns:sa-east-1:1:orders"].Message)), &msg))
	assert.JSONEq(t, `{"id":"1"}`, msg.Default)

	require.NoError(t, json.Unmarshal([]byte(aws.ToString(byTopic["arn:aws:sns:sa-east-1:1:payments"].Message)), &msg))
	var whole map[string]any
	require.NoError(t, json.Unmarshal([]byte(msg.Default), &whole))
	assert.Equal(t, "2", whole["id"])
	assert.Equal(t, "order.created", whole["type"])
}

func TestPublish_AFailingTopicFailsOnlyItsEvent(t *testing.T) {
	f := &fakeSNS{failing: map[string]bool{"down": true}}
	events := []*v2.Event{event("1", "up", nil), event("2", "down", nil), event("3", "up", nil)}

	out, err := New(f).Publish(context.Background(), events)
	require.NoError(t, err)
	assert.NoError(t, out[0].Error)
	assert.ErrorContains(t, out[1].Error, "topic down is down")
	assert.NoError(t, out[2].Error)
}

func TestPublish_NoEvents(t *testing.T) {
	f := &fakeSNS{}
	out, err := New(f).Publish(context.Background(), nil)
	require.NoError(t, err)
	assert.Empty(t, out)
	assert.Empty(t, f.inputs)
}
