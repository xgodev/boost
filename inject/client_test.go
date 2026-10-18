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

const (
	restyPath = "github.com/go-resty/resty/v2"
	redisPath = "github.com/redis/go-redis/v9"
	mongoPath = "github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v2"
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
