package mongodb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mongofact "github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// TestNew keeps the stock in the collection of Options, on the conn inject
// builds: the app closes it when it stops.
func TestNew(t *testing.T) {
	// Connect does not reach the server: no mongo is needed
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://127.0.0.1:1"))
	require.NoError(t, err)
	conn := &mongofact.Conn{Client: client, Database: client.Database("shop")}

	repo := New(conn, &Options{Collection: "units"})

	assert.Equal(t, "units", repo.(*datastore).coll.Name())
}
