package bigtable

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"cloud.google.com/go/bigtable"
	"cloud.google.com/go/bigtable/bttest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/config"
	koanf "github.com/xgodev/boost/wrapper/config/contrib/knadh/koanf/v1"
	"google.golang.org/grpc"
)

const (
	project  = "proj"
	instance = "inst"
	table    = "orders"
	family   = "cf"
)

// emulator starts an in-process Bigtable (bttest), points the SDK at it
// through BIGTABLE_EMULATOR_HOST, as an app does locally, and creates the
// test table.
func emulator(t *testing.T) {
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
	require.NoError(t, admin.CreateColumnFamily(ctx, table, family))
}

// writeAndRead writes one cell in key and asserts that reading the row
// returns it.
func writeAndRead(t *testing.T, c *bigtable.Client, key string) {
	t.Helper()
	ctx := context.Background()
	tbl := c.Open(table)

	mut := bigtable.NewMutation()
	mut.Set(family, "status", bigtable.Now(), []byte("paid"))
	require.NoError(t, tbl.Apply(ctx, key, mut))

	row, err := tbl.ReadRow(ctx, key)
	require.NoError(t, err)
	require.Len(t, row[family], 1)
	assert.Equal(t, family+":status", row[family][0].Column)
	assert.Equal(t, []byte("paid"), row[family][0].Value)
}

// loadWith loads config with the env key set to value and, on cleanup,
// reloads it with def: boost's koanf keeps keys across Load calls, so the
// value would otherwise leak into the next test.
func loadWith(t *testing.T, key, value, def string) {
	t.Helper()
	t.Setenv(key, value)
	config.Set(koanf.New())
	config.Load()
	t.Cleanup(func() {
		os.Setenv(key, def)
		config.Load()
		os.Unsetenv(key)
	})
}

func TestNewClientWithOptions_WritesAndReadsThroughEmulator(t *testing.T) {
	emulator(t)

	c, err := NewClientWithOptions(context.Background(), &Options{Project: project, Instance: instance})
	require.NoError(t, err)
	defer c.Close()

	writeAndRead(t, c, "order#1")
}

func TestNewClient_FromConfig(t *testing.T) {
	emulator(t)
	loadWith(t, "BOOST_FACTORY_BIGTABLE_PROJECT", project, "")
	loadWith(t, "BOOST_FACTORY_BIGTABLE_INSTANCE", instance, "")
	ctx := context.Background()

	c, err := NewClient(ctx)
	require.NoError(t, err)
	defer c.Close()
	writeAndRead(t, c, "order#1")

	named, err := NewClientWithConfigPath(ctx, root)
	require.NoError(t, err)
	defer named.Close()
	writeAndRead(t, named, "order#2")
}

// Without project or instance the SDK would only fail on its first RPC to
// Google; the factory refuses before dialing.
func TestNewClient_MissingProjectOrInstanceIsNotValid(t *testing.T) {
	config.Set(koanf.New())
	config.Load()
	ctx := context.Background()

	_, err := NewClient(ctx)
	assert.True(t, errors.IsNotValid(err), "got %v", err)
	_, err = NewClientWithConfigPath(ctx, root)
	assert.True(t, errors.IsNotValid(err), "got %v", err)

	for _, o := range []*Options{{Project: project}, {Instance: instance}} {
		_, err = NewClientWithOptions(ctx, o)
		assert.True(t, errors.IsNotValid(err), "got %v", err)
	}
}

func TestNewClientWithOptions_BadCredentialsFile(t *testing.T) {
	t.Setenv("BIGTABLE_EMULATOR_HOST", "") // the emulator skips credentials
	dir := t.TempDir()
	malformed := filepath.Join(dir, "malformed.json")
	require.NoError(t, os.WriteFile(malformed, []byte("not json"), 0o600))
	missing := filepath.Join(dir, "missing.json")

	for name, file := range map[string]string{"malformed": malformed, "missing": missing} {
		t.Run(name, func(t *testing.T) {
			o := &Options{Project: project, Instance: instance}
			o.Credentials.File = file
			c, err := NewClientWithOptions(context.Background(), o)

			require.Error(t, err)
			assert.False(t, errors.IsNotValid(err), "the file reached the SDK: %v", err)
			assert.Nil(t, c)
		})
	}

	o := &Options{Project: project, Instance: instance}
	o.Credentials.File = missing
	_, err := NewClientWithOptions(context.Background(), o)
	assert.ErrorContains(t, err, missing)
}

func TestNewClientWithOptions_RunsPlugins(t *testing.T) {
	emulator(t)
	ctx := context.WithValue(context.Background(), ctxKey{}, "app")

	var got []any
	plugin := func(ctx context.Context) ([]grpc.DialOption, []grpc.CallOption) {
		got = append(got, ctx.Value(ctxKey{}))
		return nil, nil
	}
	c, err := NewClientWithOptions(ctx, &Options{Project: project, Instance: instance}, plugin, plugin)
	require.NoError(t, err)
	defer c.Close()

	assert.Equal(t, []any{"app", "app"}, got)
}

type ctxKey struct{}
