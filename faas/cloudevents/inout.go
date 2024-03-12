package cloudevents

import (
	"context"

	"github.com/cloudevents/sdk-go/v2/event"
)

// InOut represents input and output events.
type InOut struct {
	In      *event.Event
	Out     *event.Event
	Err     error
	Context context.Context
}
