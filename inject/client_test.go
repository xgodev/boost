package inject

import (
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
