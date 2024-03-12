package repository

import (
	"context"

	v2 "github.com/cloudevents/sdk-go/v2" //nolint
)

// Event knows how to publish events.
type Event interface {
	Publish(context.Context, []*v2.Event) error
}
