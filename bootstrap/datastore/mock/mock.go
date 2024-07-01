package mock

import (
	"context"

	v2 "github.com/cloudevents/sdk-go/v2"
)

type mock struct{}

// NewMock creates a fake event repository implementation for testing purposes.
func NewMock() *mock {
	return &mock{}
}

// Publish simulates the publication of an event slice.
func (p *mock) Publish(ctx context.Context, events []*v2.Event) error {
	return nil
}
