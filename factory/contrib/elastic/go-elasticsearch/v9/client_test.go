package elasticsearch

import (
	"context"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost/model/errors"
)

func TestNewClientReadsOptionsFromConfig(t *testing.T) {
	n := newNode(t)
	loadEnv(t,
		env{"BOOST_FACTORY_ELASTICSEARCH_ADDRESSES", n.srv.URL, "http://127.0.0.1:9200"},
		env{"BOOST_FACTORY_ELASTICSEARCH_USERNAME", "elastic", ""},
		env{"BOOST_FACTORY_ELASTICSEARCH_PASSWORD", "secret", ""},
	)

	client, err := NewClient(context.Background())

	require.NoError(t, err)
	require.NotNil(t, client)
	reqs := n.received()
	require.Len(t, reqs, 1, "the client pings the node once on creation")
	user, pass, ok := reqs[0].BasicAuth()
	assert.True(t, ok)
	assert.Equal(t, "elastic", user)
	assert.Equal(t, "secret", pass)
}

func TestNewClientFailsOnUnreadableConfig(t *testing.T) {
	loadEnv(t, env{"BOOST_FACTORY_ELASTICSEARCH_MAX__RETRIES", "many", "3"})

	client, err := NewClient(context.Background())

	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestNewClientWithConfigPathReadsTheGivenPath(t *testing.T) {
	n := newNode(t)
	loadEnv(t, env{"APP_SEARCH_ADDRESSES", n.srv.URL, ""})

	client, err := NewClientWithConfigPath(context.Background(), "app.search")

	require.NoError(t, err)
	require.NotNil(t, client)
	assert.Len(t, n.received(), 1)
}

func TestNewClientWithConfigPathFailsOnUnreadablePath(t *testing.T) {
	loadEnv(t, env{"APP_BROKEN_MAX__RETRIES", "many", "3"})

	client, err := NewClientWithConfigPath(context.Background(), "app.broken")

	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestNewClientWithOptionsFailsOnAnInvalidConfig(t *testing.T) {
	o := &Options{Addresses: []string{"http://127.0.0.1:9200"}, CloudID: "deployment:ZXhhbXBsZS5jb20kYWJjJGRlZg=="}

	client, err := NewClientWithOptions(context.Background(), o)

	assert.Error(t, err, "addresses and cloud id together are rejected by the client")
	assert.Nil(t, client)
}

func TestNewClientWithOptionsFailsWhenTheNodeIsDown(t *testing.T) {
	n := newNode(t)
	o := n.options()
	n.srv.Close()

	client, err := NewClientWithOptions(context.Background(), o)

	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestNewClientWithOptionsFailsWhenThePingIsRejected(t *testing.T) {
	n := newNode(t, http.StatusUnauthorized)

	client, err := NewClientWithOptions(context.Background(), n.options())

	require.Error(t, err)
	assert.True(t, errors.IsServiceUnavailable(err), "a rejected ping is a service unavailable, got %v", err)
	assert.Contains(t, err.Error(), "401")
	assert.Nil(t, client)
}

func TestNewClientWithOptionsRetriesThePingOnTheConfiguredStatus(t *testing.T) {
	n := newNode(t, http.StatusServiceUnavailable, http.StatusOK)
	o := &Options{
		Addresses:     []string{n.srv.URL},
		RetryOnStatus: []int{http.StatusServiceUnavailable},
		MaxRetries:    1,
		RetryBackoff:  time.Millisecond,
	}

	client, err := NewClientWithOptions(context.Background(), o)

	require.NoError(t, err)
	require.NotNil(t, client)
	assert.Len(t, n.received(), 2)
}

func TestNewClientWithOptionsTrustsTheConfiguredCACert(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-Elastic-Product", "Elasticsearch")
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)
	ca := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: srv.Certificate().Raw})

	untrusted, err := NewClientWithOptions(context.Background(), &Options{Addresses: []string{srv.URL}, DisableRetry: true})
	require.Error(t, err, "a node with an unknown CA is refused")
	assert.Nil(t, untrusted)

	client, err := NewClientWithOptions(context.Background(), &Options{Addresses: []string{srv.URL}, DisableRetry: true, CACert: string(ca)})
	require.NoError(t, err)
	assert.NotNil(t, client)
}

func TestNewClientWithOptionsRunsPluginsInOrderOnTheBuiltClient(t *testing.T) {
	var calls []string
	var got *elasticsearch.Client
	plugin := func(name string) Plugin {
		return func(_ context.Context, c *elasticsearch.Client) error {
			calls = append(calls, name)
			got = c
			return nil
		}
	}

	client, err := NewClientWithOptions(context.Background(), newNode(t).options(), plugin("a"), plugin("b"))

	require.NoError(t, err)
	assert.Equal(t, []string{"a", "b"}, calls)
	assert.Same(t, client, got)
}

func TestNewClientWithOptionsReturnsThePluginError(t *testing.T) {
	boom := errors.New("boom")
	var after bool
	failing := func(context.Context, *elasticsearch.Client) error { return boom }
	next := func(context.Context, *elasticsearch.Client) error { after = true; return nil }

	client, err := NewClientWithOptions(context.Background(), newNode(t).options(), failing, next)

	assert.ErrorIs(t, err, boom)
	assert.Nil(t, client)
	assert.False(t, after, "plugins after the failing one must not run")
}

func TestClientConfigMapsEveryOption(t *testing.T) {
	o := &Options{
		Addresses:             []string{"http://a:9200", "http://b:9200"},
		Username:              "user",
		Password:              "pass",
		CloudID:               "cloud",
		APIKey:                "key",
		CACert:                "pem",
		RetryOnStatus:         []int{429},
		DisableRetry:          true,
		MaxRetries:            7,
		DiscoverNodesOnStart:  true,
		DiscoverNodesInterval: time.Minute,
		EnableMetrics:         true,
		EnableDebugLogger:     true,
		RetryBackoff:          10 * time.Millisecond,
	}

	cfg := clientConfig(o)

	assert.Equal(t, o.Addresses, cfg.Addresses)
	assert.Equal(t, o.Username, cfg.Username)
	assert.Equal(t, o.Password, cfg.Password)
	assert.Equal(t, o.CloudID, cfg.CloudID)
	assert.Equal(t, o.APIKey, cfg.APIKey)
	assert.Equal(t, []byte(o.CACert), cfg.CACert)
	assert.Equal(t, o.RetryOnStatus, cfg.RetryOnStatus)
	assert.Equal(t, o.DisableRetry, cfg.DisableRetry)
	assert.Equal(t, o.MaxRetries, cfg.MaxRetries)
	assert.Equal(t, o.DiscoverNodesOnStart, cfg.DiscoverNodesOnStart)
	assert.Equal(t, o.DiscoverNodesInterval, cfg.DiscoverNodesInterval)
	assert.Equal(t, o.EnableMetrics, cfg.EnableMetrics)
	assert.Equal(t, o.EnableDebugLogger, cfg.EnableDebugLogger)
	assert.Equal(t, 30*time.Millisecond, cfg.RetryBackoff(3), "the backoff grows linearly with the attempt")
	assert.IsType(t, &Logger{}, cfg.Logger)
}

func TestClientConfigLeavesCACertUnsetWhenEmpty(t *testing.T) {
	assert.Nil(t, clientConfig(&Options{}).CACert)
}
