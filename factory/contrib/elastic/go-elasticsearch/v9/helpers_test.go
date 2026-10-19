package elasticsearch

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/xgodev/boost/wrapper/config"
	koanf "github.com/xgodev/boost/wrapper/config/contrib/knadh/koanf/v1"
	"github.com/xgodev/boost/wrapper/log"
)

// node is a fake Elasticsearch node. It answers each request with the next
// status of its script (the last one repeats) and records the requests.
type node struct {
	mu       sync.Mutex
	statuses []int
	requests []*http.Request
	srv      *httptest.Server
}

// newNode starts a node answering with statuses, 200 when none is given.
func newNode(t *testing.T, statuses ...int) *node {
	t.Helper()
	if len(statuses) == 0 {
		statuses = []int{http.StatusOK}
	}
	n := &node{statuses: statuses}
	n.srv = httptest.NewServer(http.HandlerFunc(n.serve))
	t.Cleanup(n.srv.Close)
	return n
}

func (n *node) serve(w http.ResponseWriter, r *http.Request) {
	n.mu.Lock()
	status := n.statuses[min(len(n.requests), len(n.statuses)-1)]
	n.requests = append(n.requests, r)
	n.mu.Unlock()

	w.Header().Set("X-Elastic-Product", "Elasticsearch")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if strings.HasSuffix(r.URL.Path, "/_bulk") && status < http.StatusBadRequest {
		_, _ = fmt.Fprint(w, bulkResponse(r)) // intentional: the client reports a short write itself
	}
}

// bulkResponse acknowledges every document of a bulk request, one action
// line and one source line each, as the indexer pairs response items with
// the documents it sent.
func bulkResponse(r *http.Request) string {
	body, _ := io.ReadAll(r.Body) // intentional: a short body only acknowledges fewer documents
	docs := bytes.Count(body, []byte("\n")) / 2
	items := strings.TrimSuffix(strings.Repeat(`{"index":{"status":201}},`, docs), ",")
	return `{"took":1,"errors":false,"items":[` + items + `]}`
}

// received returns the requests the node got, in order.
func (n *node) received() []*http.Request {
	n.mu.Lock()
	defer n.mu.Unlock()
	return append([]*http.Request(nil), n.requests...)
}

// options points the client at the node with retries off.
func (n *node) options() *Options {
	return &Options{Addresses: []string{n.srv.URL}, DisableRetry: true}
}

// env is one environment variable a test sets; restore is the value config
// is reloaded with on cleanup, since boost's koanf keeps keys across Load
// calls and the test value would leak into the next test.
type env struct{ key, value, restore string }

// loadEnv sets the variables and loads config from them.
func loadEnv(t *testing.T, vars ...env) {
	t.Helper()
	for _, v := range vars {
		t.Setenv(v.key, v.value)
	}
	config.Set(koanf.New())
	config.Load()
	t.Cleanup(func() {
		for _, v := range vars {
			_ = os.Setenv(v.key, v.restore) // intentional: t.Setenv already proved the key is settable
		}
		config.Load()
	})
}

// entry is one line a recorder logged.
type entry struct {
	level  string
	msg    string
	fields map[string]interface{}
}

// recorder is a logger that keeps what it is asked to log.
type recorder struct {
	log.Noop
	mu      *sync.Mutex
	entries *[]entry
	fields  map[string]interface{}
}

// record installs a recorder as the boost logger for the test.
func record(t *testing.T) recorder {
	t.Helper()
	r := recorder{mu: new(sync.Mutex), entries: new([]entry)}
	log.Set(r)
	t.Cleanup(func() { log.Set(log.Noop{}) })
	return r
}

func (r recorder) add(level, msg string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	*r.entries = append(*r.entries, entry{level: level, msg: msg, fields: r.fields})
}

// lines returns the entries logged at level.
func (r recorder) lines(level string) []entry {
	r.mu.Lock()
	defer r.mu.Unlock()
	var got []entry
	for _, e := range *r.entries {
		if e.level == level {
			got = append(got, e)
		}
	}
	return got
}

func (r recorder) FromContext(context.Context) log.Logger { return r }
func (r recorder) WithTypeOf(interface{}) log.Logger      { return r }
func (r recorder) WithFields(f map[string]interface{}) log.Logger {
	r.fields = f
	return r
}
func (r recorder) Debug(args ...interface{}) { r.add("debug", fmt.Sprint(args...)) }
func (r recorder) Debugf(format string, args ...interface{}) {
	r.add("debug", fmt.Sprintf(format, args...))
}
func (r recorder) Infof(format string, args ...interface{}) {
	r.add("info", fmt.Sprintf(format, args...))
}
func (r recorder) Warn(args ...interface{})  { r.add("warn", fmt.Sprint(args...)) }
func (r recorder) Error(args ...interface{}) { r.add("error", fmt.Sprint(args...)) }
func (r recorder) Errorf(format string, args ...interface{}) {
	r.add("error", fmt.Sprintf(format, args...))
}
