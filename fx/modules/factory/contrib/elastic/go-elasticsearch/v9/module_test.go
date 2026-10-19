package elastic

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	e "github.com/xgodev/boost/factory/contrib/elastic/go-elasticsearch/v9"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

// The modules are once-guarded per process, so their options are captured
// once here and reused by every run (-count>1 included). Module comes first:
// BulkModule includes it and would otherwise take its only options.
var (
	module = Module()
	bulk   = BulkModule()
)

func TestModulesBuildTheClientAndTheBulkIndexer(t *testing.T) {
	var bulks atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Elastic-Product", "Elasticsearch")
		w.Header().Set("Content-Type", "application/json")
		if strings.HasSuffix(r.URL.Path, "/_bulk") {
			bulks.Add(1)
			_, _ = w.Write([]byte(`{"took":1,"errors":false,"items":[]}`)) // intentional: the client reports a short write itself
		}
	}))
	defer srv.Close()
	t.Setenv("BOOST_FACTORY_ELASTICSEARCH_ADDRESSES", srv.URL)
	boost.Start()

	var (
		client  *elasticsearch.Client
		indexer esutil.BulkIndexer
		plugged *elasticsearch.Client
	)
	plugin := e.Plugin(func(_ context.Context, c *elasticsearch.Client) error { plugged = c; return nil })
	app := fxtest.New(t, module, bulk, fx.Supply([]e.Plugin{plugin}), fx.Populate(&client, &indexer)).RequireStart()
	defer app.RequireStop()

	require.NotNil(t, client)
	assert.Same(t, client, plugged, "the module hands the provided plugins to the factory")
	require.NoError(t, indexer.Add(context.Background(), esutil.BulkIndexerItem{Action: "index", Body: strings.NewReader(`{}`)}))
	require.NoError(t, indexer.Close(context.Background()))
	assert.Equal(t, int32(1), bulks.Load(), "the indexer flushes through the module's client")
}
