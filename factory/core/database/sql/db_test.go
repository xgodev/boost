package sql

import (
	"context"
	"database/sql/driver"
	"errors"
	"testing"
)

var errConnect = errors.New("connection refused")

// failingConnector fails every connection, like a database that is down.
type failingConnector struct{}

func (failingConnector) Connect(context.Context) (driver.Conn, error) { return nil, errConnect }
func (failingConnector) Driver() driver.Driver                        { return failingDriver{} }

type failingDriver struct{}

func (failingDriver) Open(string) (driver.Conn, error) { return nil, errConnect }

func TestNewDBWithOptionsPingsByDefault(t *testing.T) {
	db, err := NewDBWithOptions(context.Background(), failingConnector{}, &Options{})
	if !errors.Is(err, errConnect) {
		if db != nil {
			db.Close()
		}
		t.Fatalf("NewDBWithOptions() error = %v, want %v from the ping", err, errConnect)
	}
}

func TestNewDBWithOptionsSkipPing(t *testing.T) {
	db, err := NewDBWithOptions(context.Background(), failingConnector{}, &Options{SkipPing: true})
	if err != nil {
		t.Fatalf("NewDBWithOptions() error = %v, want nil with SkipPing", err)
	}
	defer db.Close()

	if err := db.PingContext(context.Background()); !errors.Is(err, errConnect) {
		t.Fatalf("PingContext() error = %v, want %v", err, errConnect)
	}
}
