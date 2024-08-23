//go:generate mockery --name Driver --case underscore
package publisher

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type Driver interface {
	Publish(context.Context, []*cloudevents.Event) error
}
