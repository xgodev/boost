package kinesis

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/wrapper/publisher"
	"go.uber.org/fx"
)

// isolate keeps the developer's AWS files and profile out of the test.
func isolate(t *testing.T) {
	t.Helper()
	empty := filepath.Join(t.TempDir(), "empty")
	require.NoError(t, os.WriteFile(empty, nil, 0o600))
	t.Setenv("AWS_CONFIG_FILE", empty)
	t.Setenv("AWS_SHARED_CREDENTIALS_FILE", empty)
	t.Setenv("AWS_PROFILE", "")
	t.Setenv("BOOST_FACTORY_AWS_DEFAULT__REGION", "sa-east-1")
}

func TestModule_DisabledProvidesNothing(t *testing.T) {
	isolate(t)
	key := "BOOST_WRAPPER_PUBLISHER_DRIVER_KINESIS_ENABLED"
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
	isolate(t)
	boost.Start()

	var d publisher.Driver
	app := fx.New(fx.NopLogger, Module(), fx.Populate(&d))
	require.NoError(t, app.Err())
	assert.NotNil(t, d)
}
