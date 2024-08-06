package function

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type Handler[T any] func(context.Context, cloudevents.Event) (T, error)
