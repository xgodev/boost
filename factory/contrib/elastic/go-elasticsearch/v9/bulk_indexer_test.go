package elasticsearch

import (
	"context"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/elastic/go-elasticsearch/v9/esutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// flushOne adds one document to bi and closes it, which flushes it to the node.
func flushOne(t *testing.T, bi esutil.BulkIndexer) {
	t.Helper()
	require.NoError(t, bi.Add(context.Background(), esutil.BulkIndexerItem{Action: "index", Body: strings.NewReader(`{"sku":"a"}`)}))
	require.NoError(t, bi.Close(context.Background()))
}

// lastPath is the path of the last request the node got.
func lastPath(t *testing.T, n *node) string {
	t.Helper()
	reqs := n.received()
	require.NotEmpty(t, reqs)
	return reqs[len(reqs)-1].URL.Path
}

func TestNewBulkIndexerFlushesToTheConfiguredIndex(t *testing.T) {
	n := newNode(t)
	loadEnv(t, env{"BOOST_FACTORY_ELASTICSEARCH_BULK_INDEX", "orders", "index"})
	client, err := NewClientWithOptions(context.Background(), n.options())
	require.NoError(t, err)

	bi, err := NewBulkIndexer(context.Background(), client)
	require.NoError(t, err)
	flushOne(t, bi)

	assert.Equal(t, "/orders/_bulk", lastPath(t, n))
	assert.Equal(t, uint64(1), bi.Stats().NumFlushed)
}

func TestNewBulkIndexerFailsOnUnreadableConfig(t *testing.T) {
	loadEnv(t, env{"BOOST_FACTORY_ELASTICSEARCH_BULK_NUM__WORKERS", "many", strconv.Itoa(runtime.NumCPU())})

	bi, err := NewBulkIndexer(context.Background(), nil)

	assert.Error(t, err)
	assert.Nil(t, bi)
}

func TestNewBulkIndexerWithConfigPathFlushesToThePathIndex(t *testing.T) {
	n := newNode(t)
	loadEnv(t, env{"APP_SEARCH_BULK_INDEX", "events", ""})
	client, err := NewClientWithOptions(context.Background(), n.options())
	require.NoError(t, err)

	bi, err := NewBulkIndexerWithConfigPath(context.Background(), "app.search", client)
	require.NoError(t, err)
	flushOne(t, bi)

	assert.Equal(t, "/events/_bulk", lastPath(t, n))
}

func TestNewBulkIndexerWithConfigPathFailsOnUnreadablePath(t *testing.T) {
	loadEnv(t, env{"APP_BROKEN_BULK_NUM__WORKERS", "many", "1"})

	bi, err := NewBulkIndexerWithConfigPath(context.Background(), "app.broken", nil)

	assert.Error(t, err)
	assert.Nil(t, bi)
}

func TestNewBulkOptionsWithPathKeepsTheRootDefaultsThePathLeavesOut(t *testing.T) {
	loadEnv(t, env{"APP_SEARCH_BULK_INDEX", "events", ""})

	o, err := NewBulkOptionsWithPath("app.search")

	require.NoError(t, err)
	assert.Equal(t, "events", o.Index)
	assert.Equal(t, 5242880, o.FlushBytes, "flushBytes comes from boost.factory.elasticsearch.bulk")
}

func TestBulkIndexerLogsFlushErrorsAndDebugLines(t *testing.T) {
	n := newNode(t, http.StatusOK, http.StatusInternalServerError)
	client, err := NewClientWithOptions(context.Background(), n.options())
	require.NoError(t, err)
	rec := record(t)

	bi, err := NewBulkIndexerWithOptions(context.Background(), &BulkOptions{Index: "orders", NumWorkers: 1}, client)
	require.NoError(t, err)
	flushOne(t, bi)

	var flushErrs []entry
	for _, e := range rec.lines("error") {
		if strings.HasPrefix(e.msg, "ERROR: flush") {
			flushErrs = append(flushErrs, e)
		}
	}
	require.Len(t, flushErrs, 1, "a failed flush is reported once through OnError")
	assert.Contains(t, flushErrs[0].msg, "500")
	assert.NotEmpty(t, rec.lines("debug"), "the indexer's debug lines go to the boost logger")
}
