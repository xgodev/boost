package buntdb

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/buntdb"
	"github.com/xgodev/boost"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestModuleProvidesDBAndClosesItOnStop(t *testing.T) {
	boost.Start()

	var db *buntdb.DB
	app := fxtest.New(t, Module(), fx.Populate(&db))
	app.RequireStart()

	require.NotNil(t, db)
	require.NoError(t, db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set("k", "v", nil)
		return err
	}))

	app.RequireStop()

	require.ErrorIs(t, db.View(func(tx *buntdb.Tx) error { return nil }), buntdb.ErrDatabaseClosed)
}
