package memdb

import (
	"testing"

	"github.com/hashicorp/go-memdb"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

type item struct{ ID string }

func TestModuleProvidesDBFromSchemaInGraph(t *testing.T) {
	boost.Start()

	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"item": {
				Name: "item",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {Name: "id", Unique: true, Indexer: &memdb.StringFieldIndex{Field: "ID"}},
				},
			},
		},
	}

	var db *memdb.MemDB
	app := fxtest.New(t, Module(), fx.Supply(schema), fx.Populate(&db))
	app.RequireStart()
	defer app.RequireStop()

	txn := db.Txn(true)
	require.NoError(t, txn.Insert("item", &item{ID: "a"}))
	txn.Commit()

	got, err := db.Txn(false).First("item", "id", "a")
	require.NoError(t, err)
	require.Equal(t, "a", got.(*item).ID)
}
