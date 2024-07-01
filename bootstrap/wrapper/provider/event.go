package provider

import (
	"context"
	"reflect"

	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/xgodev/boost/bootstrap/repository"
)

// EventWrapperProvider represents an event wrapper provider that
// adds extra information and New Relic segment.
type EventWrapperProvider struct {
	events repository.Event
	pkg    string
	impl   string
}

// NewEventWrapperProvider creates a new event wrapper provider.
func NewEventWrapperProvider(events repository.Event) *EventWrapperProvider {
	impl := reflect.TypeOf(events).Kind().String()
	pkg := reflect.TypeOf(events).PkgPath()

	return &EventWrapperProvider{events: events, impl: impl, pkg: pkg}
}

// Publish publishes the events slice with extra information and New Relic segmentation.
func (e *EventWrapperProvider) Publish(ctx context.Context, events []*v2.Event) error {
	txn := newrelic.FromContext(ctx)

	seg := &newrelic.MessageProducerSegment{
		StartTime:            txn.StartSegmentNow(),
		Library:              e.impl,
		DestinationType:      "",
		DestinationName:      "",
		DestinationTemporary: false,
	}
	defer seg.End()

	return e.events.Publish(ctx, events)
}
