package mongo

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// TestConn_Close disconnects the client of the connection: a module stops
// it with Close when the app stops.
func TestConn_Close(t *testing.T) {
	// Connect does not reach the server: no mongo is needed
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://127.0.0.1:1"))
	if err != nil {
		t.Fatal(err)
	}
	conn := &Conn{Client: client}

	if err := conn.Close(context.Background()); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if err := client.Ping(context.Background(), nil); err != mongo.ErrClientDisconnected {
		t.Fatalf("Ping after Close = %v, want %v", err, mongo.ErrClientDisconnected)
	}
}
