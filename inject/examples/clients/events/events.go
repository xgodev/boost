// Package events publishes events over the "events" NATS connection.
package events

import "github.com/nats-io/nats.go"

// Publisher publishes events to NATS subjects.
type Publisher struct {
	conn *nats.Conn
}

// NewPublisher provides the publisher on the events connection, which
// inject builds with the health and log plugins.
// @Provide(index=0)
// @Inject(index=0, name=events, plugins=[health, log])
func NewPublisher(conn *nats.Conn) *Publisher {
	return &Publisher{conn: conn}
}

// Publish sends data to subject and waits until the server has it.
func (p *Publisher) Publish(subject string, data []byte) error {
	if err := p.conn.Publish(subject, data); err != nil {
		return err
	}
	return p.conn.Flush()
}

// Closed reports whether the connection was closed.
func (p *Publisher) Closed() bool {
	return p.conn.IsClosed()
}
