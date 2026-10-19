package freecache

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
)

func TestMain(m *testing.M) {
	os.Setenv("CONF", "../../../../../configs/config.yaml")
	boost.Start()
	os.Exit(m.Run())
}

func TestGetCopiesAFarNoteToTheNearLayer(t *testing.T) {
	ctx := context.Background()
	repo, err := New(ctx)
	require.NoError(t, err)
	d := repo.(*datastore)

	// a note the far layer alone has, like one another instance saved
	require.NoError(t, d.far.Set([]byte("shopping"), []byte("milk"), 0))
	require.Zero(t, d.near.EntryCount())

	n, err := d.Get(ctx, "shopping")
	require.NoError(t, err)
	assert.Equal(t, "milk", n.Text)
	assert.EqualValues(t, 1, d.near.EntryCount())
}
