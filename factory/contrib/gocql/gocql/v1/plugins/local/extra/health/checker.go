package health

import (
	"context"

	"github.com/gocql/gocql"
)

// Checker represents a health checker for gocql cassandra client.
type Checker struct {
	session *gocql.Session
}

// Check checks if cassandra cluster is up: it reads the node of the session
// from system.local, which every cassandra and scylla node answers.
func (c *Checker) Check(ctx context.Context) error {
	return c.session.Query("SELECT now() FROM system.local").WithContext(ctx).Exec()
}

// NewChecker returns a new health checker.
func NewChecker(session *gocql.Session) *Checker {
	return &Checker{session: session}
}
