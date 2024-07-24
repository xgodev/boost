package noop

import (
	"context"
	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/xgodev/boost/wrapper/log"
	"github.com/xgodev/boost/wrapper/publisher"
)

// client represents a noop client.
type client struct {
}

// New creates a new noop client.
func New() publisher.Driver {
	return &client{}
}

// Publish publishes an event slice.
func (p *client) Publish(ctx context.Context, outs []*v2.Event) (err error) {
	logger := log.FromContext(ctx)
	logger.Debugf("published on the noop")
	return nil
}
