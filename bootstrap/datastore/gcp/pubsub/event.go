package pubsub

import (
	"cloud.google.com/go/pubsub"
	"github.com/xgodev/boost/bootstrap/repository"
)

// NewEvent returns a initialized client
func NewEvent(c *pubsub.Client) repository.Event {
	return NewClient(c)
}
