package inject

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectEntries_KeepsInjectAnnotations(t *testing.T) {
	root := t.TempDir()
	files := map[string]string{
		"go.mod": "module example.com/collect\n\ngo 1.21\n",
		"stock/memory/memory.go": `package memory

// New is the default implementation.
// @Provide(index=0)
// @Default
func New() any { return nil }

// Ignored carries no inject annotation.
// @Other(key=value)
func Ignored() {}
`,
	}
	for name, body := range files {
		p := filepath.Join(root, name)
		require.NoError(t, os.MkdirAll(filepath.Dir(p), 0o755))
		require.NoError(t, os.WriteFile(p, []byte(body), 0o644))
	}

	entries, err := CollectEntries(root)
	require.NoError(t, err)

	var names []string
	for _, entry := range entries {
		for _, ann := range entry.Annotations {
			names = append(names, ann.Name+" "+gid(entry))
		}
	}
	sort.Strings(names)
	assert.Equal(t, []string{
		"Default example.com/collect/stock/memory.New",
		"Provide example.com/collect/stock/memory.New",
	}, names)
}
