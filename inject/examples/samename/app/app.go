// Package app consumes the publishers of two packages named boost.
package app

import (
	changed "github.com/xgodev/boost/inject/examples/samename/publisher/pricechanged/boost"
	received "github.com/xgodev/boost/inject/examples/samename/publisher/pricereceived/boost"
)

// Subjects holds the subject each publisher sends to.
type Subjects struct {
	Received string
	Changed  string
}

// NewSubjects reads the subject of each publisher.
// @Provide(index=0)
func NewSubjects(r *received.Publisher, c *changed.Publisher) *Subjects {
	return &Subjects{Received: r.Subject, Changed: c.Subject}
}
