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
func (p *client) Publish(ctx context.Context, outs []*v2.Event) (res []publisher.PublishOutput, err error) {
	logger := log.FromContext(ctx).WithTypeOf(*p)
	for _, out := range outs {
		logger.WithField("event", out).Debugf("published event to noop topic %s", out.Subject())
		res = append(res, publisher.PublishOutput{Event: out})
	}
	logger.Debugf("published all on the noop")
	return res, nil
}
