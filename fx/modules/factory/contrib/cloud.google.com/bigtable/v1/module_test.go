package bigtable

import (
	"context"
	"os"
	"testing"

	"cloud.google.com/go/bigtable"
	"cloud.google.com/go/bigtable/bttest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost/wrapper/config"
	koanf "github.com/xgodev/boost/wrapper/config/contrib/knadh/koanf/v1"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// leaf holds the module's only options: Module() is once-guarded, so a
// second call in this binary would be empty.
var leaf = Module()

const (
	project  = "proj"
	instance = "inst"
	table    = "orders"
)

// setup starts an in-process Bigtable (bttest) with the test table, points
// the SDK at it through BIGTABLE_EMULATOR_HOST and loads project and
// instance into boost.factory.bigtable.
func setup(t *testing.T) {
	t.Helper()
	srv, err := bttest.NewServer("localhost:0")
	require.NoError(t, err)
	t.Cleanup(srv.Close)
	t.Setenv("BIGTABLE_EMULATOR_HOST", srv.Addr)

	ctx := context.Background()
	admin, err := bigtable.NewAdminClient(ctx, project, instance)
	require.NoError(t, err)
	defer admin.Close()
	require.NoError(t, admin.CreateTable(ctx, table))

	t.Setenv("BOOST_FACTORY_BIGTABLE_PROJECT", project)
	t.Setenv("BOOST_FACTORY_BIGTABLE_INSTANCE", instance)
	config.Set(koanf.New())
	config.Load()
	t.Cleanup(func() {
		// boost's koanf keeps keys across Load calls: put the defaults back.
		os.Setenv("BOOST_FACTORY_BIGTABLE_PROJECT", "")
		os.Setenv("BOOST_FACTORY_BIGTABLE_INSTANCE", "")
		config.Load()
	})
}

func TestModuleGraphIsComplete(t *testing.T) {
	var c *bigtable.Client
	require.NoError(t, fx.ValidateApp(leaf, fx.Populate(&c)))
}

func TestModuleClosesClientOnStop(t *testing.T) {
	setup(t)

	var c *bigtable.Client
	app := fxtest.New(t, leaf, fx.Populate(&c))
	app.RequireStart()

	ctx := context.Background()
	_, err := c.Open(table).ReadRow(ctx, "order#1")
	require.NoError(t, err, "the client works while the app runs")

	app.RequireStop()

	_, err = c.Open(table).ReadRow(ctx, "order#1")
	assert.Equal(t, codes.Canceled, status.Code(err), "the client is closed after stop: %v", err)
}
