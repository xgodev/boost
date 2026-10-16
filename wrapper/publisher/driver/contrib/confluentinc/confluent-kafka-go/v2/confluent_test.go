package confluent

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/stretchr/testify/suite"
	"github.com/xgodev/boost/wrapper/publisher"
)

type ConfluentSuite struct {
	suite.Suite
}

func TestConfluentSuite(t *testing.T) {
	suite.Run(t, new(ConfluentSuite))
}

// queued is a message the fake accepted, waiting for its delivery report.
type queued struct {
	index int
	msg   *kafka.Message
	ch    chan kafka.Event
}

// fakeProducer stands in for *kafka.Producer: Produce only enqueues, and the
// delivery report of every queued message arrives later, from another
// goroutine and in reverse order, the way librdkafka reports across
// partitions. Indexes are Produce call order.
type fakeProducer struct {
	mu sync.Mutex
	// expected is how many messages are queued before reports are sent.
	expected int
	// syncErrs fails Produce itself for a call index.
	syncErrs map[int]error
	// deliveryErrs sets TopicPartition.Error on the report of a call index.
	deliveryErrs map[int]error
	// silent never sends a report (broker gone after enqueue).
	silent bool
	// onAllQueued runs once every expected message is queued.
	onAllQueued func()

	calls           int
	queue           []queued
	nilDeliveryChan bool
}

func (f *fakeProducer) Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	index := f.calls
	f.calls++
	if err := f.syncErrs[index]; err != nil {
		return err
	}
	if deliveryChan == nil {
		f.nilDeliveryChan = true
	}
	f.queue = append(f.queue, queued{index: index, msg: msg, ch: deliveryChan})
	if len(f.queue) != f.expected {
		return nil
	}
	if f.onAllQueued != nil {
		f.onAllQueued()
	}
	if !f.silent && !f.nilDeliveryChan {
		go f.report(append([]queued(nil), f.queue...))
	}
	return nil
}

func (f *fakeProducer) report(queue []queued) {
	for i := len(queue) - 1; i >= 0; i-- {
		q := queue[i]
		// librdkafka hands back a new message carrying the original Opaque.
		report := &kafka.Message{
			TopicPartition: q.msg.TopicPartition,
			Key:            q.msg.Key,
			Value:          q.msg.Value,
			Headers:        q.msg.Headers,
			Opaque:         q.msg.Opaque,
		}
		report.TopicPartition.Error = f.deliveryErrs[q.index]
		q.ch <- report
	}
}

func (f *fakeProducer) usedNilDeliveryChan() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.nilDeliveryChan
}

func newEvent(id string) *v2.Event {
	e := v2.NewEvent()
	e.SetID(id)
	e.SetSource("test")
	e.SetType("test")
	e.SetSubject("orders")
	_ = e.SetData(v2.ApplicationJSON, map[string]string{"id": id}) // intentional: a map always marshals
	return &e
}

// testCtx bounds a test that would otherwise hang on a missing report.
func (s *ConfluentSuite) testCtx() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	s.T().Cleanup(cancel)
	return ctx
}

func (s *ConfluentSuite) assertEventsInOrder(events []*v2.Event, outputs []publisher.PublishOutput) {
	s.Require().Len(outputs, len(events))
	for i := range events {
		s.Same(events[i], outputs[i].Event, "output %d", i)
	}
}

func (s *ConfluentSuite) TestPublishPutsDeliveryFailureOnTheMatchingEvent() {
	errTimedOut := kafka.NewError(kafka.ErrMsgTimedOut, "Local: Message timed out", false)
	events := []*v2.Event{newEvent("evt-0"), newEvent("evt-timed-out"), newEvent("evt-2")}
	fake := &fakeProducer{expected: len(events), deliveryErrs: map[int]error{1: errTimedOut}}

	outputs, err := newDriver(fake, &Options{}).Publish(s.testCtx(), events)

	s.Require().NoError(err)
	s.False(fake.usedNilDeliveryChan(), "driver must pass a delivery channel to Produce")
	s.assertEventsInOrder(events, outputs)
	s.NoError(outputs[0].Error)
	var kerr kafka.Error
	s.Require().True(errors.As(outputs[1].Error, &kerr), "want the delivery report error, got %v", outputs[1].Error)
	s.Equal(kafka.ErrMsgTimedOut, kerr.Code())
	s.NoError(outputs[2].Error)
}

func (s *ConfluentSuite) TestPublishSucceedsWhenEveryMessageIsDelivered() {
	events := []*v2.Event{newEvent("evt-0"), newEvent("evt-1")}
	fake := &fakeProducer{expected: len(events)}

	outputs, err := newDriver(fake, &Options{}).Publish(s.testCtx(), events)

	s.Require().NoError(err)
	s.False(fake.usedNilDeliveryChan(), "driver must pass a delivery channel to Produce")
	s.assertEventsInOrder(events, outputs)
	for i, out := range outputs {
		s.NoError(out.Error, "output %d", i)
	}
}

func (s *ConfluentSuite) TestPublishKeepsSynchronousProduceErrorAndWaitsForTheRest() {
	errQueueFull := kafka.NewError(kafka.ErrQueueFull, "Local: Queue full", false)
	errUnknownPartition := kafka.NewError(kafka.ErrUnknownPartition, "Local: Unknown partition", false)
	events := []*v2.Event{newEvent("evt-queue-full"), newEvent("evt-1"), newEvent("evt-bad-partition")}
	fake := &fakeProducer{
		expected:     2,
		syncErrs:     map[int]error{0: errQueueFull},
		deliveryErrs: map[int]error{2: errUnknownPartition},
	}

	outputs, err := newDriver(fake, &Options{}).Publish(s.testCtx(), events)

	s.Require().NoError(err)
	s.assertEventsInOrder(events, outputs)
	var kerr kafka.Error
	s.Require().True(errors.As(outputs[0].Error, &kerr), "want the Produce error, got %v", outputs[0].Error)
	s.Equal(kafka.ErrQueueFull, kerr.Code())
	s.NoError(outputs[1].Error)
	s.Require().True(errors.As(outputs[2].Error, &kerr), "want the delivery report error, got %v", outputs[2].Error)
	s.Equal(kafka.ErrUnknownPartition, kerr.Code())
}

func (s *ConfluentSuite) TestPublishFailsPendingEventsWhenContextEndsBeforeDelivery() {
	ctx, cancel := context.WithCancel(s.testCtx())
	events := []*v2.Event{newEvent("evt-0"), newEvent("evt-1")}
	fake := &fakeProducer{expected: len(events), silent: true, onAllQueued: cancel}

	outputs, err := newDriver(fake, &Options{}).Publish(ctx, events)

	s.Require().NoError(err)
	s.assertEventsInOrder(events, outputs)
	for i, out := range outputs {
		s.ErrorIs(out.Error, context.Canceled, "output %d", i)
	}
}
