package inject

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost/annotation"
	"github.com/xgodev/boost/model/errors"
)

const (
	restyPath  = "github.com/go-resty/resty/v2"
	redisPath  = "github.com/redis/go-redis/v9"
	mongoPath  = "github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v2"
	natsPath   = "github.com/nats-io/nats.go"
	pubsubPath = "cloud.google.com/go/pubsub/v2"
	gocqlPath  = "github.com/gocql/gocql"
	gocql2Path = "github.com/apache/cassandra-gocql-driver/v2"
)

func TestClientFactory(t *testing.T) {
	testCases := []struct {
		name      string
		imports   map[string]string
		tp        string
		wantFound bool
		wantBuild string
	}{
		{name: "resty", imports: map[string]string{"resty": restyPath}, tp: "*resty.Client", wantFound: true,
			wantBuild: "github.com/xgodev/boost/factory/contrib/go-resty/resty/v2.NewClientWithConfigPath"},
		{name: "aliased import", imports: map[string]string{"r": restyPath}, tp: "*r.Client", wantFound: true,
			wantBuild: "github.com/xgodev/boost/factory/contrib/go-resty/resty/v2.NewClientWithConfigPath"},
		{name: "redis cluster", imports: map[string]string{"redis": redisPath}, tp: "*redis.ClusterClient", wantFound: true,
			wantBuild: "github.com/xgodev/boost/factory/contrib/redis/go-redis/v9.NewClusterClientWithConfigPath"},
		{name: "mongo v2", imports: map[string]string{"mongo": mongoPath}, tp: "*mongo.Conn", wantFound: true,
			wantBuild: mongoPath + ".NewConnWithConfigPath"},
		{name: "nats", imports: map[string]string{"nats": natsPath}, tp: "*nats.Conn", wantFound: true,
			wantBuild: "github.com/xgodev/boost/factory/contrib/nats-io/nats.go/v1.NewConnWithConfigPath"},
		{name: "pubsub v2", imports: map[string]string{"pubsub": pubsubPath}, tp: "*pubsub.Client", wantFound: true,
			wantBuild: "github.com/xgodev/boost/factory/contrib/cloud.google.com/pubsub/v2.NewClientWithConfigPath"},
		{name: "gocql", imports: map[string]string{"gocql": gocqlPath}, tp: "*gocql.Session", wantFound: true,
			wantBuild: "github.com/xgodev/boost/factory/contrib/gocql/gocql/v1.NewSessionWithConfigPath"},
		{name: "gocql v2", imports: map[string]string{"gocql": gocql2Path}, tp: "*gocql.Session", wantFound: true,
			wantBuild: "github.com/xgodev/boost/factory/contrib/apache/cassandra-gocql-driver/v2.NewSessionWithConfigPath"},
		{name: "aws sns", imports: map[string]string{"sns": "github.com/aws/aws-sdk-go-v2/service/sns"}, tp: "*sns.Client", wantFound: true,
			wantBuild: "github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/client/sns.NewSDKClientWithConfigPath"},
		{name: "aws sqs", imports: map[string]string{"sqs": "github.com/aws/aws-sdk-go-v2/service/sqs"}, tp: "*sqs.Client", wantFound: true,
			wantBuild: "github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/client/sqs.NewSDKClientWithConfigPath"},
		{name: "aws kinesis", imports: map[string]string{"kinesis": "github.com/aws/aws-sdk-go-v2/service/kinesis"}, tp: "*kinesis.Client", wantFound: true,
			wantBuild: "github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1/client/kinesis.NewSDKClientWithConfigPath"},
		{name: "not a pointer", imports: map[string]string{"resty": restyPath}, tp: "resty.Client"},
		{name: "another type of the package", imports: map[string]string{"resty": restyPath}, tp: "*resty.Request"},
		{name: "same name, another package", imports: map[string]string{"resty": "example.com/resty"}, tp: "*resty.Client"},
		{name: "local type", tp: "*Client"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			entry := at("quote", "New", []string{tc.tp}, []string{"*Repository"})
			entry.Imports = tc.imports

			f, ok := clientFactory(entry, tc.tp)
			require.Equal(t, tc.wantFound, ok)
			if ok {
				assert.Equal(t, tc.wantBuild, f.Build)
			}
		})
	}
}

func TestClientFactories_PluginNames(t *testing.T) {
	assert.Equal(t, []string{"datadog", "health", "log", "otel", "prometheus", "requestid", "retry"},
		clientFactories[restyPath+".Client"].PluginNames())
	assert.Equal(t, []string{"datadog", "health", "otel", "prometheus"},
		clientFactories[redisPath+".ClusterClient"].PluginNames())
	assert.Equal(t, []string{"health", "log"}, clientFactories[natsPath+".Conn"].PluginNames())
	assert.Equal(t, []string{"datadog", "log", "prometheus"}, clientFactories[pubsubPath+".Client"].PluginNames())
	assert.Equal(t, []string{"health"}, clientFactories[gocqlPath+".Session"].PluginNames())
	assert.Equal(t, []string{"health"}, clientFactories[gocql2Path+".Session"].PluginNames())
	assert.Equal(t, []string{"datadog"},
		clientFactories["github.com/aws/aws-sdk-go-v2/service/sqs.Client"].PluginNames())
}

func TestParsePlugins(t *testing.T) {
	testCases := []struct {
		raw     string
		want    []string
		invalid bool
	}{
		{raw: "", want: nil},
		{raw: "[otel, prometheus]", want: []string{"otel", "prometheus"}},
		{raw: "otel", want: []string{"otel"}},
		{raw: "[]", want: []string{}},
		{raw: "[otel", invalid: true},
		{raw: "otel]", invalid: true},
		{raw: "[otel,]", invalid: true},
		{raw: "[log, log]", invalid: true},
	}

	for _, tc := range testCases {
		t.Run(tc.raw, func(t *testing.T) {
			got, err := parsePlugins(tc.raw)
			if tc.invalid {
				assert.True(t, errors.IsNotValid(err), "unexpected error: %v", err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

// restyConsumer builds New in package client/<pkg>, a constructor taking a
// *resty.Client.
func restyConsumer(pkg string, anns ...annotation.Annotation) annotation.Entry {
	entry := at("client/"+pkg, "New", []string{"*resty.Client"}, []string{"*Repository"}, append(anns, provide("index=0"))...)
	entry.Imports = map[string]string{"resty": restyPath}
	return entry
}

func errorContains(t *testing.T, err error, parts ...string) {
	t.Helper()
	require.Error(t, err)
	for _, part := range parts {
		assert.True(t, strings.Contains(err.Error(), part), "%q does not contain %q", err.Error(), part)
	}
}

func TestClients_DefaultAndNamed(t *testing.T) {
	clients, err := Clients(newTestGraph(t, []annotation.Entry{
		restyConsumer("weather", inject("index=0")),
		restyConsumer("quote", inject("index=0, name=quote, plugins=[otel, prometheus]")),
	}))
	require.NoError(t, err)
	require.Len(t, clients, 2)

	def, quote := clients[0], clients[1]
	assert.Equal(t, "", def.Name)
	assert.Nil(t, def.Plugins)
	assert.Equal(t, "boost.factory.resty", def.Root())
	assert.Equal(t, "boost.factory.resty.plugins.log", def.PluginRoot("log"))
	assert.Equal(t, "DefaultClient", def.Ident())
	assert.Equal(t, "client/github.com/go-resty/resty/v2/default_client.go", def.File())
	assert.Equal(t, []string{testModule + "/client/weather.New"}, def.Consumers)

	assert.Equal(t, "quote", quote.Name)
	assert.Equal(t, []string{"otel", "prometheus"}, quote.Plugins)
	assert.Equal(t, "boost.factory.resty.clients.quote", quote.Root())
	assert.Equal(t, "boost.factory.resty.clients.quote.plugins.otel", quote.PluginRoot("otel"))
	assert.Equal(t, "QuoteClient", quote.Ident())
	assert.Equal(t, "client/github.com/go-resty/resty/v2/quote_client.go", quote.File())
	assert.Equal(t, []string{testModule + "/client/quote.New"}, quote.Consumers)

	imp, call := quote.Module("example.com/app/gen")
	assert.Equal(t, "example.com/app/gen/client/github.com/go-resty/resty/v2", imp.Path)
	assert.Equal(t, imp.Alias+".QuoteClientModule()", call)
}

func TestClients_Shared(t *testing.T) {
	clients, err := Clients(newTestGraph(t, []annotation.Entry{
		restyConsumer("quote", inject("index=0, name=quote, plugins=[log]")),
		restyConsumer("price", inject("index=0, name=quote")),
	}))
	require.NoError(t, err)
	require.Len(t, clients, 1, "every consumer of a name receives the same client")
	assert.Equal(t, []string{"log"}, clients[0].Plugins, "a consumer without plugins takes the declared ones")
	assert.Equal(t, []string{testModule + "/client/price.New", testModule + "/client/quote.New"}, clients[0].Consumers)
}

func TestClients_RedisClientAndCluster(t *testing.T) {
	entry := at("cache", "New", []string{"*redis.Client", "*redis.ClusterClient"}, []string{"*Cache"},
		provide("index=0"), inject("index=0, name=cache"), inject("index=1, name=cache, plugins=[health]"))
	entry.Imports = map[string]string{"redis": redisPath}

	clients, err := Clients(newTestGraph(t, []annotation.Entry{entry}))
	require.NoError(t, err)
	require.Len(t, clients, 2, "a client and a cluster client of one name are two clients")
	assert.Equal(t, "client/github.com/redis/go-redis/v9/cache_client.go", clients[0].File())
	assert.Equal(t, "client/github.com/redis/go-redis/v9/cache_clusterclient.go", clients[1].File())
	assert.Equal(t, "CacheClusterClient", clients[1].Ident())
	assert.Equal(t, []string{"health"}, clients[1].Plugins)
}

func TestClients_AppProviderWins(t *testing.T) {
	provider := at("client/http", "NewQuote", nil, []string{"*resty.Client"}, provide("index=0, name=quote"))
	provider.Imports = map[string]string{"resty": restyPath}
	g := newTestGraph(t, []annotation.Entry{provider, restyConsumer("quote", inject("index=0, name=quote"))})

	clients, err := Clients(g)
	require.NoError(t, err)
	assert.Empty(t, clients)
	assert.Equal(t, []string{provider.Path + ".NewQuote"}, edgesOf(g.Vertices[testModule+"/client/quote.New"].Incoming()))
}

func TestClients_NoInjectIsLeftToTheApp(t *testing.T) {
	clients, err := Clients(newTestGraph(t, []annotation.Entry{restyConsumer("quote")}))
	require.NoError(t, err)
	assert.Empty(t, clients)
}

func TestClients_Rules(t *testing.T) {
	provider := at("client/http", "NewQuote", nil, []string{"*resty.Client"}, provide("index=0, name=quote"))
	provider.Imports = map[string]string{"resty": restyPath}
	notClient := at("greet", "New", []string{"string"}, []string{"*Greeter"}, provide("index=0"), inject("index=0, plugins=[log]"))
	group := at("client/many", "New", []string{"[]*resty.Client"}, []string{"*Many"}, provide("index=0"), inject("index=0, group=clients, plugins=[log]"))
	group.Imports = map[string]string{"resty": restyPath}
	provided := at("client/http", "NewQuote", nil, []string{"*resty.Client"}, provide("index=0, plugins=[log]"))
	provided.Imports = map[string]string{"resty": restyPath}

	testCases := []struct {
		name    string
		entries []annotation.Entry
		parts   []string
	}{
		{name: "plugins on a parameter that is not a client", entries: []annotation.Entry{notClient},
			parts: []string{"plugins", "string", "github.com/go-resty/resty/v2.Client"}},
		{name: "plugins on a client the app provides", entries: []annotation.Entry{provider, restyConsumer("quote", inject("index=0, name=quote, plugins=[log]"))},
			parts: []string{"provided by", provider.Path + ".NewQuote"}},
		{name: "group with plugins", entries: []annotation.Entry{group},
			parts: []string{"group", "plugins"}},
		{name: "name an env var cannot reach", entries: []annotation.Entry{restyConsumer("quote", inject("index=0, name=Quote"))},
			parts: []string{"Quote"}},
		{name: "unknown plugin", entries: []annotation.Entry{restyConsumer("quote", inject("index=0, plugins=[nope]"))},
			parts: []string{"nope", "datadog, health, log, otel, prometheus, requestid, retry"}},
		{name: "unbalanced list", entries: []annotation.Entry{restyConsumer("quote", inject("index=0, plugins=[otel"))},
			parts: []string{"[otel"}},
		{name: "different lists", entries: []annotation.Entry{
			restyConsumer("quote", inject("index=0, name=quote, plugins=[log]")),
			restyConsumer("price", inject("index=0, name=quote, plugins=[otel]")),
		}, parts: []string{testModule + "/client/price.New", testModule + "/client/quote.New", "[log]", "[otel]"}},
		{name: "declared empty against a list", entries: []annotation.Entry{
			restyConsumer("quote", inject("index=0, name=quote, plugins=[log]")),
			restyConsumer("price", inject("index=0, name=quote, plugins=[]")),
		}, parts: []string{"[log]", "[]"}},
		{name: "plugins on @Provide", entries: []annotation.Entry{provided},
			parts: []string{"@Provide", "plugins"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewGraphFromEntries(context.Background(), tc.entries)
			assert.True(t, errors.IsNotValid(err), "unexpected error: %v", err)
			errorContains(t, err, tc.parts...)
		})
	}
}

func TestClients_DefaultNameCollides(t *testing.T) {
	_, err := NewGraphFromEntries(context.Background(), []annotation.Entry{
		restyConsumer("weather", inject("index=0")),
		restyConsumer("quote", inject("index=0, name=default")),
	})
	assert.True(t, errors.IsAlreadyExists(err), "unexpected error: %v", err)
	errorContains(t, err, "default_client.go", testModule+"/client/quote.New", testModule+"/client/weather.New")
}

func TestGenerator_RenderClient(t *testing.T) {
	weather := restyConsumer("weather", inject("index=0"))
	quote := restyConsumer("quote", inject("index=0, name=quote, plugins=[otel, prometheus]"))

	files, err := NewGenerator(newTestGraph(t, []annotation.Entry{weather, quote}), testOptions("/mod")).Render(context.Background())
	require.NoError(t, err)

	named := files["client/github.com/go-resty/resty/v2/quote_client.go"]
	require.NotNil(t, named, "the named client is rendered")
	golden(t, "testdata/client/quote_client.go.golden", named)

	def := string(files["client/github.com/go-resty/resty/v2/default_client.go"])
	assert.Contains(t, def, `defaultClientRoot = "boost.factory.resty"`)
	assert.Contains(t, def, "fx.Provide(newDefaultClient)", "the default client has no fx name")
	assert.Contains(t, def, "factory.NewClientWithConfigPath(ctx, defaultClientRoot)", "no plugin is wired when none is listed")

	clients, err := Clients(newTestGraph(t, []annotation.Entry{weather, quote}))
	require.NoError(t, err)
	imp, call := clients[1].Module("example.com/app/gen")
	consumer := string(files["client_/quote/new_module.go"])
	assert.Contains(t, consumer, imp.Alias+` "example.com/app/gen/client/github.com/go-resty/resty/v2"`)
	assert.Contains(t, consumer, call+",", "the consumer includes the client it injects")
	assert.Contains(t, consumer, "fx.ParamTags(`name:\"quote\"`)")
	_, defCall := clients[0].Module("example.com/app/gen")
	assert.NotContains(t, consumer, defCall, "a consumer includes only its own clients")

	root := string(files["module.go"])
	assert.Equal(t, 1, strings.Count(root, call+","), "the root includes the client once")
	assert.Equal(t, 1, strings.Count(root, defCall+","))
}

// TestGenerator_RenderClosingClients renders a named client of each factory
// whose client holds a connection: its module closes it when the app stops.
func TestGenerator_RenderClosingClients(t *testing.T) {
	testCases := []struct {
		name, pkg, path, tp, inject, file string
	}{
		{name: "nats", pkg: "nats", path: natsPath, tp: "*nats.Conn", inject: "index=0, name=events, plugins=[health, log]",
			file: "client/github.com/nats-io/nats.go/events_conn.go"},
		{name: "pubsub", pkg: "pubsub", path: pubsubPath, tp: "*pubsub.Client", inject: "index=0, name=orders, plugins=[log, prometheus]",
			file: "client/cloud.google.com/go/pubsub/v2/orders_client.go"},
		{name: "gocql", pkg: "gocql", path: gocqlPath, tp: "*gocql.Session", inject: "index=0, name=store, plugins=[health]",
			file: "client/github.com/gocql/gocql/store_session.go"},
		{name: "gocql v2", pkg: "gocql", path: gocql2Path, tp: "*gocql.Session", inject: "index=0, name=ledger, plugins=[health]",
			file: "client/github.com/apache/cassandra-gocql-driver/v2/ledger_session.go"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			entry := at("client/"+tc.pkg, "New", []string{tc.tp}, []string{"*Repository"}, provide("index=0"), inject(tc.inject))
			entry.Imports = map[string]string{tc.pkg: tc.path}

			files, err := NewGenerator(newTestGraph(t, []annotation.Entry{entry}), testOptions("/mod")).Render(context.Background())
			require.NoError(t, err)

			src := files[tc.file]
			require.NotNil(t, src, "the client is rendered at %s", tc.file)
			golden(t, "testdata/client/"+filepath.Base(tc.file)+".golden", src)
		})
	}
}

// TestClientRegistryCompiles renders every client of the registry with all
// its plugins and builds it, so a wrong constructor name or plugin type in
// the registry fails here instead of in an application.
func TestClientRegistryCompiles(t *testing.T) {
	if testing.Short() {
		t.Skip("builds every factory and plugin")
	}

	tmpl, err := NewTemplate()
	require.NoError(t, err)
	dir, err := os.MkdirTemp(".", "zz_clientcheck_")
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.RemoveAll(dir) })

	keys := make([]string, 0, len(clientFactories))
	for key := range clientFactories {
		keys = append(keys, key)
	}
	for i, key := range keys {
		f := clientFactories[key]
		for _, c := range []Client{{Factory: f}, {Factory: f, Name: "check", Plugins: f.PluginNames()}} {
			src, err := ExecuteTemplate(tmpl, "client", clientData(c))
			require.NoError(t, err, "%s", c)
			file := filepath.Join(dir, strings.ReplaceAll(key, "/", "_"), filepath.Base(c.File()))
			require.NoError(t, os.MkdirAll(filepath.Dir(file), os.ModePerm), "%d", i)
			require.NoError(t, os.WriteFile(file, src, 0o644))
		}
	}

	out, err := exec.Command("go", "build", "./"+filepath.Base(dir)+"/...").CombinedOutput()
	require.NoError(t, err, "%s", out)
}
