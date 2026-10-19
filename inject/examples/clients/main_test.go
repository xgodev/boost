package main

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"cloud.google.com/go/pubsub/pstest"
	natsserver "github.com/nats-io/nats-server/v2/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/inject"
	"github.com/xgodev/boost/inject/examples/clients/events"
	"github.com/xgodev/boost/inject/examples/clients/gen"
	"github.com/xgodev/boost/inject/examples/clients/ledger"
	"github.com/xgodev/boost/inject/examples/clients/orders"
	"github.com/xgodev/boost/inject/examples/clients/store"
	"go.uber.org/fx"
)

// TestGeneratedModulesStartAndStopTheClients boots the app against an
// embedded NATS server, the Pub/Sub fake and a fake CQL node, reached by
// both gocql drivers: each client works after start and is closed after
// stop.
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
	boost.Start() // loads the config keys the generated modules registered

	ctx := context.Background()
	var publisher *events.Publisher
	var topics *orders.Topics
	var cluster *store.Cluster
	var ledgerCluster *ledger.Cluster
	app := fx.New(gen.Module(), fx.Populate(&publisher, &topics, &cluster, &ledgerCluster), fx.NopLogger)
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

	require.NoError(t, app.Stop(ctx))

	assert.True(t, publisher.Closed(), "the NATS connection is closed on stop")
	assert.Error(t, topics.Create(ctx, "after-stop"), "the Pub/Sub client is closed on stop")
	assert.True(t, cluster.Closed(), "the gocql session is closed on stop")
	assert.True(t, ledgerCluster.Closed(), "the gocql v2 session is closed on stop")
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
