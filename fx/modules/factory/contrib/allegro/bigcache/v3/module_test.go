package bigcache

import (
	"testing"

	"github.com/allegro/bigcache/v3"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestModuleProvidesCache(t *testing.T) {
	boost.Start()

	var cache *bigcache.BigCache
	app := fxtest.New(t, Module(), fx.Populate(&cache))
	app.RequireStart()
	defer app.RequireStop()

	require.NotNil(t, cache)
	require.NoError(t, cache.Set("k", []byte("v")))
	got, err := cache.Get("k")
	require.NoError(t, err)
	require.Equal(t, []byte("v"), got)
}
