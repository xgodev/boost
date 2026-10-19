package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"cloud.google.com/go/pubsub/pstest"
	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	natsserver "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/inject"
	"github.com/xgodev/boost/inject/examples/publisher/domain/model/entity"
	"github.com/xgodev/boost/inject/examples/publisher/domain/publisher/ordercreated"
	"github.com/xgodev/boost/inject/examples/publisher/gen"
	ordercreatedboost "github.com/xgodev/boost/inject/examples/publisher/publisher/ordercreated/boost"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/publisher"
	"go.uber.org/fx"
)

// TestMain boots boost: the driver key and the clients' keys are read from
// config.
func TestMain(m *testing.M) {
	boost.Start()
	os.Exit(m.Run())
}

const (
	// driverEnv chooses the publisher's driver (boost.wrapper.publisher.driver.name).
	driverEnv     = "BOOST_WRAPPER_PUBLISHER_DRIVER_NAME"
	natsURLEnv    = "BOOST_FACTORY_NATS_URL"
	projectEnv    = "BOOST_FACTORY_GCP_PUBSUB_API__OPTIONS_PROJECT__ID"
	emulatorEnv   = "PUBSUB_EMULATOR_HOST"
	childEnv      = "INJECT_PUBLISHER_EXAMPLE_CHILD"
	project       = "example"
	driverPackage = "github.com/xgodev/boost/wrapper/publisher/driver/"
)

// TestPublisherDriver boots the example in a new process per value of the
// driver key: config is read once per process. The nats and pubsub drivers
// publish to a NATS server and a Pub/Sub fake this test runs.
func TestPublisherDriver(t *testing.T) {
	opts := natsserver.DefaultTestOptions
	opts.Port = -1
	srv := natsserver.RunServer(&opts)
	t.Cleanup(srv.Shutdown)

	fake := pstest.NewServer()
	t.Cleanup(func() { _ = fake.Close() })

	for _, value := range []string{"", "noop", "nats", "pubsub", "kafka", "nope"} {
		t.Run("driver="+value, func(t *testing.T) {
			cmd := exec.Command(os.Args[0], "-test.run=^TestPublisherDriverChild$", "-test.v")
			cmd.Env = append(os.Environ(), childEnv+"=1", natsURLEnv+"="+srv.ClientURL(), emulatorEnv+"="+fake.Addr, projectEnv+"="+project)
			if value != "" {
				cmd.Env = append(cmd.Env, driverEnv+"="+value)
			}
			out, err := cmd.CombinedOutput()
			require.NoError(t, err, "%s", out)
			assert.Contains(t, string(out), "--- PASS: TestPublisherDriverChild")
		})
	}
}

// TestPublisherWithoutCgo builds the example with CGO_ENABLED=0: every
// driver but kafka is built in, and choosing kafka fails the boot saying
// it needs cgo.
func TestPublisherWithoutCgo(t *testing.T) {
	if testing.Short() {
		t.Skip("builds the example without cgo")
	}

	cmd := exec.Command("go", "test", "-count=1", "-run=^TestPublisherDriverChild$", "-v", ".")
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0", childEnv+"=1", driverEnv+"=kafka")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "%s", out)
	assert.Contains(t, string(out), "--- PASS: TestPublisherDriverChild")
	assert.Contains(t, string(out), "without cgo")
}

// TestPublisherDriverChild runs in the process TestPublisherDriver starts,
// with the driver key set or unset.
func TestPublisherDriverChild(t *testing.T) {
	if os.Getenv(childEnv) != "1" {
		t.Skip("run by TestPublisherDriver")
	}

	var (
		orders ordercreated.Publisher
		pub    *publisher.Publisher
	)
	app := fx.New(gen.Module(), fx.NopLogger, fx.Populate(&orders, &pub))
	ctx := context.Background()
	order := entity.Order{ID: "1", Item: "book"}

	switch value := os.Getenv(driverEnv); value {
	case "", "noop":
		require.NoError(t, app.Err())
		assert.Equal(t, driverPackage+"extra/noop", driverOf(t, pub), "noop is the default")
		require.NoError(t, orders.Publish(ctx, order))

	case "nats":
		require.NoError(t, app.Err())
		assert.Equal(t, driverPackage+"contrib/nats-io/nats.go/v1", driverOf(t, pub))

		conn, err := nats.Connect(os.Getenv(natsURLEnv))
		require.NoError(t, err)
		t.Cleanup(conn.Close)
		sub, err := conn.SubscribeSync(ordercreatedboost.Subject)
		require.NoError(t, err)
		require.NoError(t, conn.Flush())

		require.NoError(t, orders.Publish(ctx, order))
		msg, err := sub.NextMsg(5 * time.Second)
		require.NoError(t, err, "the event reaches NATS")
		assert.Contains(t, string(msg.Data), `"item":"book"`)

	case "pubsub":
		require.NoError(t, app.Err())
		assert.Equal(t, driverPackage+"contrib/cloud.google.com/pubsub/v2", driverOf(t, pub), "v2 is the default Pub/Sub version")

		subscriber := pubsubSubscriber(t, ctx)
		require.NoError(t, orders.Publish(ctx, order))
		assert.Contains(t, string(receive(t, ctx, subscriber)), `"item":"book"`, "the event reaches Pub/Sub")

	case "kafka":
		if !cgoEnabled {
			require.Error(t, app.Err(), "kafka needs cgo")
			assert.True(t, errors.IsNotSupported(app.Err()), "got %v", app.Err())
			assert.Contains(t, app.Err().Error(), "boost.wrapper.publisher.driver.name")
			t.Log("kafka fails the boot without cgo")
			return
		}
		require.NoError(t, app.Err(), "the producer is created without a reachable broker")
		assert.Equal(t, driverPackage+"contrib/confluentinc/confluent-kafka-go/v2", driverOf(t, pub))

	default:
		require.Error(t, app.Err(), "an unknown driver fails the boot")
		assert.True(t, errors.IsNotValid(app.Err()), "got %v", app.Err())
		assert.Contains(t, app.Err().Error(), "boost.wrapper.publisher.driver.name")
		assert.Contains(t, app.Err().Error(), `"`+value+`"`)
		assert.Contains(t, app.Err().Error(), "kafka, nats, noop, pubsub")
	}
}

// driverOf returns the package of the driver the publisher was built with.
func driverOf(t *testing.T, p *publisher.Publisher) string {
	t.Helper()
	require.NotNil(t, p)
	driver := reflect.ValueOf(p).Elem().FieldByName("driver")
	require.False(t, driver.IsNil())
	return driver.Elem().Type().Elem().PkgPath()
}

// pubsubSubscriber creates the event's topic and a subscription to it on
// the Pub/Sub fake.
func pubsubSubscriber(t *testing.T, ctx context.Context) *pubsub.Subscriber {
	t.Helper()
	client, err := pubsub.NewClient(ctx, project)
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })

	topic := "projects/" + project + "/topics/" + ordercreatedboost.Subject
	_, err = client.TopicAdminClient.CreateTopic(ctx, &pubsubpb.Topic{Name: topic})
	require.NoError(t, err)
	sub := "projects/" + project + "/subscriptions/orders"
	_, err = client.SubscriptionAdminClient.CreateSubscription(ctx, &pubsubpb.Subscription{Name: sub, Topic: topic})
	require.NoError(t, err)
	return client.Subscriber(sub)
}

// receive returns the data of the first message of the subscription.
func receive(t *testing.T, ctx context.Context, sub *pubsub.Subscriber) []byte {
	t.Helper()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var data []byte
	err := sub.Receive(ctx, func(_ context.Context, m *pubsub.Message) {
		data = m.Data
		m.Ack()
		cancel()
	})
	require.NoError(t, err)
	require.NotNil(t, data, "no message within 10s")
	return data
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

	const hint = "gen/ is stale: run go generate ./inject/examples/publisher"
	for name, src := range files {
		assert.Equal(t, string(src), onDisk[name], "%s: %s", name, hint)
		delete(onDisk, name)
	}
	assert.Empty(t, onDisk, "files not produced by the generator: %s", hint)
}
