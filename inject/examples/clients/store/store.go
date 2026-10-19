// Package store reads the Cassandra/ScyllaDB cluster through the "store"
// session.
package store

import "github.com/gocql/gocql"

// Cluster queries the cluster.
type Cluster struct {
	session *gocql.Session
}

// NewCluster provides the cluster on the store session, which inject
// builds with the health plugin.
// @Provide(index=0)
// @Inject(index=0, name=store, plugins=[health])
func NewCluster(session *gocql.Session) *Cluster {
	return &Cluster{session: session}
}

// Version returns the release version of the node the session reaches.
func (c *Cluster) Version() (string, error) {
	var version string
	err := c.session.Query("SELECT release_version FROM system.local WHERE key='local'").Scan(&version)
	return version, err
}

// Closed reports whether the session was closed.
func (c *Cluster) Closed() bool {
	return c.session.Closed()
}
