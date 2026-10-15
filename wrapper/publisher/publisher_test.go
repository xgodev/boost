package publisher_test

import (
	"context"
	"errors"
	"testing"

	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/xgodev/boost/extra/middleware"
	boosterrors "github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/publisher"
	"github.com/xgodev/boost/wrapper/publisher/mocks"
)

type PublisherSuite struct {
	suite.Suite
}

func TestPublisherSuite(t *testing.T) {
	suite.Run(t, new(PublisherSuite))
}

func newEvent(id string) *v2.Event {
	e := v2.NewEvent()
	e.SetID(id)
	e.SetSource("test")
	e.SetType("test")
	return &e
}

// recorder is a middleware that keeps what the publish chain returned.
type recorder struct {
	outputs []publisher.PublishOutput
	err     error
}

func (m *recorder) Exec(ctx *middleware.AnyErrorContext[[]publisher.PublishOutput],
	exec middleware.AnyErrorExecFunc[[]publisher.PublishOutput],
	returnFunc middleware.AnyErrorReturnFunc[[]publisher.PublishOutput]) ([]publisher.PublishOutput, error) {
	m.outputs, m.err = ctx.Next(exec, returnFunc)
	return m.outputs, m.err
}

func (s *PublisherSuite) TestPublishReturnsEveryPerEventError() {
	errQueueFull := errors.New("queue full")
	errConversion := errors.New("conversion failed")
	failedA, ok, failedB := newEvent("evt-queue-full"), newEvent("evt-delivered"), newEvent("evt-bad-payload")
	events := []*v2.Event{failedA, ok, failedB}

	drv := mocks.NewDriver(s.T())
	drv.On("Publish", mock.Anything, events).Return([]publisher.PublishOutput{
		{Event: failedA, Error: errQueueFull},
		{Event: ok},
		{Event: failedB, Error: errConversion},
	}, nil)

	err := publisher.New(drv).Publish(context.Background(), events)

	s.Require().Error(err)
	s.True(boosterrors.IsInternal(err))
	s.ErrorIs(err, errQueueFull)
	s.ErrorIs(err, errConversion)
	s.Contains(err.Error(), "evt-queue-full")
	s.Contains(err.Error(), "evt-bad-payload")
	s.NotContains(err.Error(), "evt-delivered")
}

// Middlewares (e.g. prometheus) count each event from the outputs and skip
// counting when the chain fails, so the per-event failure must reach them as
// outputs, not as a chain error.
func (s *PublisherSuite) TestPublishKeepsPerEventOutputsForMiddleware() {
	errQueueFull := errors.New("queue full")
	failed, ok := newEvent("evt-queue-full"), newEvent("evt-delivered")
	events := []*v2.Event{failed, ok}
	outputs := []publisher.PublishOutput{
		{Event: failed, Error: errQueueFull},
		{Event: ok},
	}

	drv := mocks.NewDriver(s.T())
	drv.On("Publish", mock.Anything, events).Return(outputs, nil)

	rec := &recorder{}
	err := publisher.New(drv, rec).Publish(context.Background(), events)

	s.ErrorIs(err, errQueueFull)
	s.NoError(rec.err)
	s.Equal(outputs, rec.outputs)
}

func (s *PublisherSuite) TestPublishSucceedsWhenEveryEventIsDelivered() {
	a, b := newEvent("evt-a"), newEvent("evt-b")
	events := []*v2.Event{a, b}

	drv := mocks.NewDriver(s.T())
	drv.On("Publish", mock.Anything, events).Return([]publisher.PublishOutput{
		{Event: a},
		{Event: b},
	}, nil)

	s.NoError(publisher.New(drv).Publish(context.Background(), events))
}

func (s *PublisherSuite) TestPublishReturnsDriverError() {
	errBrokerDown := errors.New("broker down")
	events := []*v2.Event{newEvent("evt-a")}

	drv := mocks.NewDriver(s.T())
	drv.On("Publish", mock.Anything, events).Return(nil, errBrokerDown)

	s.ErrorIs(publisher.New(drv).Publish(context.Background(), events), errBrokerDown)
}
