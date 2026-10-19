package health

import (
	"context"

	"github.com/apache/cassandra-gocql-driver/v2"
)

// query reads one row any node has, so it answers while the node serves CQL.
const query = "SELECT release_version FROM system.local"

// Checker reports whether a gocql session reaches Cassandra.
type Checker struct {
	session *gocql.Session
}

// Check runs a read on the session; an error means Cassandra did not answer.
func (c *Checker) Check(ctx context.Context) error {
	return c.session.Query(query).ExecContext(ctx)
}

// NewChecker returns a checker for session.
func NewChecker(session *gocql.Session) *Checker {
	return &Checker{session: session}
}
