package nats

import (
	"context"
	"github.com/nats-io/nats.go"

	"github.com/xgodev/boost/bootstrap/repository"
)

// NewEvent returns an initialized NATS client that implements event repository.
func NewEvent(ctx context.Context, conn *nats.Conn) repository.Event {
	return NewClient(conn)
}
