// Package tcp is a health checker of its own: it dials a tcp dependency.
package tcp

import (
	"context"
	"net"
	"time"

	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/config"
)

// Checker dials the dependency; it implements health.Checker.
type Checker struct {
	addr string
}

// New provides the checker of the dependency at app.dependency.addr.
// @Provide(index=0)
func New() *Checker {
	return &Checker{addr: config.String(addr)}
}

// Check fails when the dependency does not accept a connection in a second.
func (c *Checker) Check(ctx context.Context) error {
	var d net.Dialer
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	conn, err := d.DialContext(ctx, "tcp", c.addr)
	if err != nil {
		return errors.Wrap(err, errors.ServiceUnavailablef("dependency %s", c.addr))
	}
	return conn.Close()
}
