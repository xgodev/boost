package elasticsearch

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeServer answers the client's ping like an Elasticsearch node would.
func fakeServer(t *testing.T) *Options {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-Elastic-Product", "Elasticsearch")
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)
	return &Options{Addresses: []string{srv.URL}, DisableRetry: true}
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

	client, err := NewClientWithOptions(context.Background(), fakeServer(t), plugin("a"), plugin("b"))

	require.NoError(t, err)
	assert.Equal(t, []string{"a", "b"}, calls)
	assert.Same(t, client, got)
}

func TestNewClientWithOptionsPanicsOnPluginError(t *testing.T) {
	boom := errors.New("boom")
	var after bool
	failing := func(context.Context, *elasticsearch.Client) error { return boom }
	next := func(context.Context, *elasticsearch.Client) error { after = true; return nil }
	opts := fakeServer(t)

	assert.PanicsWithError(t, boom.Error(), func() {
		_, _ = NewClientWithOptions(context.Background(), opts, failing, next)
	})
	assert.False(t, after, "plugins after the failing one must not run")
}
