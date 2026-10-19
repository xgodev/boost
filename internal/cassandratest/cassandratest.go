// Package cassandratest gives the tests of the gocql v2 factory, its plugins
// and its fx module a Cassandra to talk to.
package cassandratest

import (
	"net"
	"os"
	"strconv"
	"testing"
)

// Env names the variable holding the host:port of a Cassandra the
// integration tests may use, e.g. 127.0.0.1:9042 for
// `docker run -p 9042:9042 cassandra:5.0`.
const Env = "BOOST_TEST_CASSANDRA"

// Addr returns the host and port in Env, and skips the test, saying how to
// run it, when Env is unset.
func Addr(t *testing.T) (string, int) {
	t.Helper()
	v := os.Getenv(Env)
	if v == "" {
		t.Skipf("%s is unset: start a Cassandra (docker run -p 9042:9042 cassandra:5.0) and export %s=127.0.0.1:9042", Env, Env)
	}
	host, p, err := net.SplitHostPort(v)
	if err != nil {
		t.Fatalf("%s=%q is not host:port: %v", Env, v, err)
	}
	port, err := strconv.Atoi(p)
	if err != nil {
		t.Fatalf("%s=%q has a non-numeric port: %v", Env, v, err)
	}
	return host, port
}

// ClosedPort returns a local port nothing listens on, so a dial to it is
// refused at once.
func ClosedPort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	if err := l.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	return port
}
