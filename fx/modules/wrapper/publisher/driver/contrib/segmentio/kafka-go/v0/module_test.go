package kafka

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/wrapper/publisher"
	"go.uber.org/fx"
)

func TestModule_DisabledProvidesNothing(t *testing.T) {
	key := "BOOST_WRAPPER_PUBLISHER_DRIVER_KAFKAGO_ENABLED"
	t.Setenv(key, "false")
	boost.Start()
	// boost keeps config keys across loads: reload the default, or "false"
	// would leak into the next test.
	t.Cleanup(func() {
		os.Setenv(key, "true")
		boost.Start()
	})

	var d publisher.Driver
	app := fx.New(fx.NopLogger, Module(), fx.Populate(&d))
	assert.Error(t, app.Err())
}

func TestModule_ProvidesTheDriver(t *testing.T) {
	boost.Start()

	var d publisher.Driver
	app := fx.New(fx.NopLogger, Module(), fx.Populate(&d))
	require.NoError(t, app.Err())
	assert.NotNil(t, d)
}
