package main

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"cloud.google.com/go/bigtable"
	"cloud.google.com/go/bigtable/bttest"
	"cloud.google.com/go/pubsub/pstest"
	natsserver "github.com/nats-io/nats-server/v2/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/inject"
	"github.com/xgodev/boost/inject/examples/clients/catalog"
	"github.com/xgodev/boost/inject/examples/clients/events"
	"github.com/xgodev/boost/inject/examples/clients/gen"
	"github.com/xgodev/boost/inject/examples/clients/ledger"
	"github.com/xgodev/boost/inject/examples/clients/orders"
	"github.com/xgodev/boost/inject/examples/clients/store"
	"go.uber.org/fx"
)

// TestGeneratedModulesStartAndStopTheClients boots the app against an
// embedded NATS server, the Pub/Sub fake, a fake CQL node, reached by both
// gocql drivers, and the in-process Bigtable emulator: each client works
// after start and is closed after stop.
func TestGeneratedModulesStartAndStopTheClients(t *testing.T) {
	opts := natsserver.DefaultTestOptions
	opts.Port = -1
	nats := natsserver.RunServer(&opts)
	defer nats.Shutdown()

	ps := pstest.NewServer()
	defer ps.Close()

	t.Setenv("BOOST_FACTORY_NATS_CLIENTS_EVENTS_URL", nats.ClientURL())
	t.Setenv("PUBSUB_EMULATOR_HOST", ps.Addr)
	t.Setenv("BOOST_FACTORY_GCP_PUBSUB_CLIENTS_ORDERS_API__OPTIONS_PROJECT__ID", "example")
	cqlPort := strconv.Itoa(startCQLNode(t))
	t.Setenv("BOOST_FACTORY_GOCQL_CLIENTS_STORE_PORT", cqlPort)
	t.Setenv("BOOST_FACTORY_GOCQL_CLIENTS_LEDGER_PORT", cqlPort)
	startBigtable(t, "example", "catalog")
	t.Setenv("BOOST_FACTORY_BIGTABLE_CLIENTS_CATALOG_PROJECT", "example")
	t.Setenv("BOOST_FACTORY_BIGTABLE_CLIENTS_CATALOG_INSTANCE", "catalog")
	boost.Start() // loads the config keys the generated modules registered

	ctx := context.Background()
	var publisher *events.Publisher
	var topics *orders.Topics
	var cluster *store.Cluster
	var ledgerCluster *ledger.Cluster
	var prices *catalog.Prices
	app := fx.New(gen.Module(), fx.Populate(&publisher, &topics, &cluster, &ledgerCluster, &prices), fx.NopLogger)
	require.NoError(t, app.Err())
	require.NoError(t, app.Start(ctx))

	require.NoError(t, publisher.Publish("orders.created", []byte("42")))
	require.NoError(t, topics.Create(ctx, "orders"))
	version, err := cluster.Version()
	require.NoError(t, err)
	assert.Equal(t, cqlReleaseVersion, version)
	version, err = ledgerCluster.Version()
	require.NoError(t, err)
	assert.Equal(t, cqlReleaseVersion, version)
	require.NoError(t, prices.Set(ctx, "sku-1", "9.90"))
	price, err := prices.Get(ctx, "sku-1")
	require.NoError(t, err)
	assert.Equal(t, "9.90", price)

	require.NoError(t, app.Stop(ctx))

	assert.True(t, publisher.Closed(), "the NATS connection is closed on stop")
	assert.Error(t, topics.Create(ctx, "after-stop"), "the Pub/Sub client is closed on stop")
	assert.True(t, cluster.Closed(), "the gocql session is closed on stop")
	assert.True(t, ledgerCluster.Closed(), "the gocql v2 session is closed on stop")
	_, err = prices.Get(ctx, "sku-1")
	assert.Error(t, err, "the Bigtable client is closed on stop")
}

// startBigtable starts the in-process Bigtable emulator (bttest), points the
// SDK at it through BIGTABLE_EMULATOR_HOST and creates the prices table in
// project and instance.
func startBigtable(t *testing.T, project, instance string) {
	t.Helper()
	srv, err := bttest.NewServer("localhost:0")
	require.NoError(t, err)
	t.Cleanup(srv.Close)
	t.Setenv("BIGTABLE_EMULATOR_HOST", srv.Addr)

	ctx := context.Background()
	admin, err := bigtable.NewAdminClient(ctx, project, instance)
	require.NoError(t, err)
	defer func() { _ = admin.Close() }()
	require.NoError(t, admin.CreateTable(ctx, catalog.Table))
	require.NoError(t, admin.CreateColumnFamily(ctx, catalog.Table, catalog.Family))
}

// TestGeneratedCodeIsUpToDate fails when gen/ differs from what the
// generator produces for the annotations in this example. Run go generate.
func TestGeneratedCodeIsUpToDate(t *testing.T) {
	ctx := context.Background()

	modulePath, moduleDir, err := inject.FindModule(".")
	require.NoError(t, err)

	entries, err := inject.CollectEntries(".")
	require.NoError(t, err)

	g, err := inject.NewGraphFromEntries(ctx, entries)
	require.NoError(t, err)

	out, err := filepath.Abs("gen")
	require.NoError(t, err)

	files, err := inject.NewGenerator(g, inject.GeneratorOptions{
		ModulePath: modulePath,
		ModuleDir:  moduleDir,
		OutputDir:  out,
	}).Render(ctx)
	require.NoError(t, err)

	onDisk := map[string]string{}
	require.NoError(t, filepath.WalkDir(out, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, err := filepath.Rel(out, p)
		if err != nil {
			return err
		}
		content, err := os.ReadFile(p)
		onDisk[filepath.ToSlash(rel)] = string(content)
		return err
	}))

	const hint = "gen/ is stale: run go generate ./inject/examples/clients"
	for name, src := range files {
		assert.Equal(t, string(src), onDisk[name], "%s: %s", name, hint)
		delete(onDisk, name)
	}
	assert.Empty(t, onDisk, "files not produced by the generator: %s", hint)
}
