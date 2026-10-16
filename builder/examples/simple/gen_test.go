package simple

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost/builder"
	"github.com/xgodev/boost/inject"
)

// TestGeneratedCodeIsUpToDate fails when gen/ differs from what cmd/builder
// produces for the annotations in this example. Run go generate.
func TestGeneratedCodeIsUpToDate(t *testing.T) {
	ctx := context.Background()

	modulePath, moduleDir, err := inject.FindModule(".")
	require.NoError(t, err)

	entries, err := builder.CollectEntries(".")
	require.NoError(t, err)

	g, err := inject.NewGraphFromEntries(ctx, entries)
	require.NoError(t, err)

	spec, err := builder.NewSpec(entries, g)
	require.NoError(t, err)

	out, err := filepath.Abs("gen")
	require.NoError(t, err)
	opts := inject.GeneratorOptions{ModulePath: modulePath, ModuleDir: moduleDir, OutputDir: out}

	files, err := inject.NewGenerator(g, opts).Render(ctx)
	require.NoError(t, err)
	doors, err := builder.NewGenerator(spec, opts).Render(ctx)
	require.NoError(t, err)
	for name, src := range doors {
		files[name] = src
	}

	onDisk := map[string]string{}
	require.NoError(t, filepath.WalkDir(out, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, err := filepath.Rel(out, p)
		if err != nil {
			return err
		}
		content, err := os.ReadFile(p)
		onDisk[filepath.ToSlash(rel)] = string(content)
		return err
	}))

	const hint = "gen/ is stale: run go generate ./builder/examples/simple"
	for name, src := range files {
		assert.Equal(t, string(src), onDisk[name], "%s: %s", name, hint)
		delete(onDisk, name)
	}
	assert.Empty(t, onDisk, "files not produced by the generator: %s", hint)
}
