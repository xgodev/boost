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

// Consume handles events.
// @Function(adapter=nats)
func (h *Handler) Consume() {}

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
		"Function example.com/collect/app.Handler.Consume",
		"GRPC example.com/collect/app.Health",
		"Provide example.com/collect/app.New",
		"REST example.com/collect/app.Handler.Get",
	}, names)
}
