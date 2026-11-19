package pgx

import (
	"context"
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/xgodev/boost"
	"github.com/xgodev/boost/wrapper/config"
)

func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

// closedPort returns a local port nothing listens on.
func closedPort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	_ = l.Close()
	return port
}

// loadClient registers the keys of a client root, sets its env vars and
// reloads the config; the config is reloaded again after the env is restored.
func loadClient(t *testing.T, path string, env map[string]string) {
	t.Helper()
	ConfigAdd(path)
	t.Cleanup(config.Load)
	for k, v := range env {
		t.Setenv(k, v)
	}
	config.Load()
}

func TestNewDBWithConfigPathPingFalseSkipsThePing(t *testing.T) {
	loadClient(t, root+".clients.noping", map[string]string{
		"BOOST_FACTORY_PGX_CLIENTS_NOPING_CONNECTSTRING": fmt.Sprintf("postgres://127.0.0.1:%d/db?sslmode=disable&connect_timeout=2", closedPort(t)),
		"BOOST_FACTORY_PGX_CLIENTS_NOPING_PING":          "false",
	})

	db, err := NewDBWithConfigPath(context.Background(), root+".clients.noping")
	if err != nil {
		t.Fatalf("NewDBWithConfigPath() error = %v, want nil with ping=false", err)
	}
	defer db.Close()

	if err := db.PingContext(context.Background()); err == nil {
		t.Fatal("PingContext() error = nil, want connection error on a closed port")
	}
}

func TestNewDBWithConfigPathPingsByDefault(t *testing.T) {
	loadClient(t, root+".clients.ping", map[string]string{
		"BOOST_FACTORY_PGX_CLIENTS_PING_CONNECTSTRING": fmt.Sprintf("postgres://127.0.0.1:%d/db?sslmode=disable&connect_timeout=2", closedPort(t)),
	})

	db, err := NewDBWithConfigPath(context.Background(), root+".clients.ping")
	if err == nil {
		db.Close()
		t.Fatal("NewDBWithConfigPath() error = nil, want connection error from the ping on a closed port")
	}
}
