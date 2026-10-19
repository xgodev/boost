package inject

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost/annotation"
	"github.com/xgodev/boost/model/errors"
)

// natsPath and pubsubPath are declared in client_test.go.
const (
	publisherPath = "github.com/xgodev/boost/wrapper/publisher"
	kafkaPath     = "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// publisherConsumer builds New in package pkg, a constructor taking a
// *publisher.Publisher.
func publisherConsumer(pkg string, anns ...annotation.Annotation) annotation.Entry {
	entry := at(pkg, "New", []string{"*publisher.Publisher"}, []string{"*Events"}, append(anns, provide("index=0"))...)
	entry.Imports = map[string]string{"publisher": publisherPath}
	return entry
}

func TestPublisherDrivers(t *testing.T) {
	assert.Equal(t, "kafka, nats, noop, pubsub", publisherDriverNames())

	testCases := []struct {
		client     string
		root       string
		build      string
		constraint string
	}{
		{client: natsPath + ".Conn", root: "boost.factory.nats",
			build: boostFactories + "nats-io/nats.go/v1.NewConnWithConfigPath"},
		{client: pubsubPath + ".Client", root: "boost.factory.gcp.pubsub",
			build: boostFactories + "cloud.google.com/pubsub/v2.NewClientWithConfigPath"},
		{client: kafkaPath + ".Producer", root: "boost.factory.confluent",
			build: boostFactories + "confluentinc/confluent-kafka-go/v2.NewProducerWithConfigPath", constraint: "cgo"},
	}
	for _, tc := range testCases {
		t.Run(tc.client, func(t *testing.T) {
			f, ok := clientFactories[tc.client]
			require.True(t, ok, "the driver's client is one inject builds, so an @Inject consumer shares it")
			assert.Equal(t, tc.root, f.Root)
			assert.Equal(t, tc.build, f.Build)
			assert.Equal(t, tc.constraint, f.Constraint)
		})
	}

	for _, d := range publisherDrivers {
		if d.Client != "" {
			_, ok := clientFactories[d.Client]
			assert.True(t, ok, "driver %s: client %s is not registered", d.Name, d.Client)
		}
	}
}

func TestPublisherOf_Consumers(t *testing.T) {
	p := PublisherOf(newTestGraph(t, []annotation.Entry{
		publisherConsumer("order", inject("index=0")),
		publisherConsumer("billing", inject("index=0")),
	}))
	assert.Equal(t, []string{testModule + "/billing.New", testModule + "/order.New"}, p.Consumers)
	assert.Equal(t, "client/github.com/xgodev/boost/wrapper/publisher/publisher.go", p.File())

	imp, call := p.Module("example.com/app/gen")
	assert.Equal(t, "example.com/app/gen/client/github.com/xgodev/boost/wrapper/publisher", imp.Path)
	assert.Equal(t, imp.Alias+".PublisherModule()", call)
}

func TestPublisherOf_LeftToTheApp(t *testing.T) {
	provider := at("events", "NewPublisher", nil, []string{"*publisher.Publisher"}, provide("index=0"))
	provider.Imports = map[string]string{"publisher": publisherPath}

	t.Run("the app provides it", func(t *testing.T) {
		g := newTestGraph(t, []annotation.Entry{provider, publisherConsumer("order", inject("index=0"))})
		p := PublisherOf(g)
		assert.Empty(t, p.Consumers)
		assert.Equal(t, []string{provider.Path + ".NewPublisher"}, edgesOf(g.Vertices[testModule+"/order.New"].Incoming()))
	})

	t.Run("no @Inject", func(t *testing.T) {
		p := PublisherOf(newTestGraph(t, []annotation.Entry{publisherConsumer("order")}))
		assert.Empty(t, p.Consumers)
	})
}

func TestPublisherOf_Rules(t *testing.T) {
	t.Run("a named publisher is not built", func(t *testing.T) {
		_, err := NewGraphFromEntries(context.Background(), []annotation.Entry{publisherConsumer("order", inject("index=0, name=orders"))})
		assert.True(t, errors.IsNotFound(err), "unexpected error: %v", err)
	})

	t.Run("plugins", func(t *testing.T) {
		_, err := NewGraphFromEntries(context.Background(), []annotation.Entry{publisherConsumer("order", inject("index=0, plugins=[log]"))})
		assert.True(t, errors.IsNotValid(err), "unexpected error: %v", err)
		errorContains(t, err, "plugins", "*publisher.Publisher")
	})
}

func TestGenerator_RenderPublisher(t *testing.T) {
	entries := []annotation.Entry{publisherConsumer("order", inject("index=0")), publisherConsumer("billing", inject("index=0"))}
	files, err := NewGenerator(newTestGraph(t, entries), testOptions("/mod")).Render(context.Background())
	require.NoError(t, err)

	dir := "client/github.com/xgodev/boost/wrapper/publisher/"
	for _, name := range []string{"publisher.go", "publisher_kafka.go", "publisher_kafka_nocgo.go"} {
		require.NotNil(t, files[dir+name], "%s is rendered", name)
		golden(t, "testdata/publisher/"+name+".golden", files[dir+name])
	}

	producer := string(files["client/github.com/confluentinc/confluent-kafka-go/v2/kafka/default_producer.go"])
	assert.True(t, strings.HasPrefix(producer, GeneratedHeader+"\n\n//go:build cgo\n"),
		"the kafka client builds only with cgo, so an app without cgo still builds:\n%s", producer)
	conn := string(files["client/github.com/nats-io/nats.go/default_conn.go"])
	assert.True(t, strings.HasPrefix(conn, GeneratedHeader+"\n\npackage nats\n"), "a client without constraint has no build line:\n%s", conn)
	assert.NotNil(t, files["client/cloud.google.com/go/pubsub/v2/default_client.go"])

	p := PublisherOf(newTestGraph(t, entries))
	imp, call := p.Module("example.com/app/gen")
	order := string(files["order/new_module.go"])
	assert.Contains(t, order, imp.Alias+` "example.com/app/gen/client/github.com/xgodev/boost/wrapper/publisher"`)
	assert.Contains(t, order, call+",", "the consumer includes the publisher it injects")

	root := string(files["module.go"])
	assert.Equal(t, 1, strings.Count(root, call+","), "the root includes the publisher once")
	assert.NotContains(t, root, "DefaultConnModule()", "a driver's client is wired by the switch, only when chosen")
	assert.NotContains(t, root, "DefaultProducerModule()")
}

func TestGenerator_RenderPublisherSharesTheClient(t *testing.T) {
	stream := at("stream", "New", []string{"*nats.Conn"}, []string{"*Stream"}, provide("index=0"), inject("index=0"))
	stream.Imports = map[string]string{"nats": natsPath}
	entries := []annotation.Entry{publisherConsumer("order", inject("index=0")), stream}

	files, err := NewGenerator(newTestGraph(t, entries), testOptions("/mod")).Render(context.Background())
	require.NoError(t, err)

	clients, err := Clients(newTestGraph(t, entries))
	require.NoError(t, err)
	require.Len(t, clients, 1)
	_, call := clients[0].Module("example.com/app/gen")
	assert.Contains(t, string(files["client/github.com/xgodev/boost/wrapper/publisher/publisher.go"]), call,
		"the nats driver gets the default *nats.Conn the app injects")
	assert.Equal(t, 1, strings.Count(string(files["module.go"]), call+","), "the root includes the app's client once")
}
