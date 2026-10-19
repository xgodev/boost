package builder

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectEntries_KeepsDoorAndInjectAnnotations(t *testing.T) {
	root := t.TempDir()
	files := map[string]string{
		"go.mod": "module example.com/collect\n\ngo 1.21\n",
		"app/app.go": `package app

// New builds the handler.
// @Provide(index=0)
func New() *Handler { return &Handler{} }

// Handler serves the doors.
type Handler struct{}

// Get reads.
// @REST(method=GET, path=/x)
func (h *Handler) Get() error { return nil }

// Health serves gRPC.
// @GRPC(register=Register)
func Health() {}

// Check answers a tool call.
// @MCP
func (h *Handler) Check() {}

// Consume handles events.
// @Function(adapter=nats)
func (h *Handler) Consume() {}

// NewDefault is the default implementation of *Handler.
// @Provide(index=0)
// @Default
func NewDefault() *Handler { return &Handler{} }

// Options configures the handler.
// @Config
type Options struct {
	Size int
}

// Ignored carries no door.
// @Other(key=value)
func Ignored() {}
`,
	}
	for name, body := range files {
		path := filepath.Join(root, name)
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
		require.NoError(t, os.WriteFile(path, []byte(body), 0o644))
	}

	entries, err := CollectEntries(root)
	require.NoError(t, err)

	var names []string
	for _, entry := range entries {
		for _, ann := range entry.Annotations {
			names = append(names, ann.Name+" "+declName(entry))
		}
	}
	sort.Strings(names)
	assert.Equal(t, []string{
		"Config example.com/collect/app.Options",
		"Default example.com/collect/app.NewDefault",
		"Function example.com/collect/app.Handler.Consume",
		"GRPC example.com/collect/app.Health",
		"MCP example.com/collect/app.Handler.Check",
		"Provide example.com/collect/app.New",
		"Provide example.com/collect/app.NewDefault",
		"REST example.com/collect/app.Handler.Get",
	}, names)
}

func TestCollectEntries_KeepsTheAdapterListWhole(t *testing.T) {
	root := t.TempDir()
	files := map[string]string{
		"go.mod": "module example.com/collect\n\ngo 1.21\n",
		"app/app.go": `package app

// Consume handles events.
// @Function(adapter=[pubsub, nats], app=worker)
func Consume() {}
`,
	}
	for name, body := range files {
		path := filepath.Join(root, name)
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
		require.NoError(t, os.WriteFile(path, []byte(body), 0o644))
	}

	entries, err := CollectEntries(root)
	require.NoError(t, err)
	require.Len(t, entries, 1)
	require.Len(t, entries[0].Annotations, 1)

	var args functionArgs
	require.NoError(t, decode(entries[0].Annotations[0], declName(entries[0]), functionKeys, &args))
	assert.Equal(t, "[pubsub, nats]", args.Adapter, "the comma of the list does not split the annotation")
	assert.Equal(t, "worker", args.App)
}
