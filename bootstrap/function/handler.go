package function

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type Handler func(context.Context, cloudevents.Event) (*cloudevents.Event, error)
