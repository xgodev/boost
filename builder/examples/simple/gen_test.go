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

// TestGeneratedCodeIsUpToDate fails when gen/ or the mains under cmd/ differ
// from what cmd/builder produces for the annotations in this example. Run go
// generate.
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
	cmd, err := filepath.Abs("cmd")
	require.NoError(t, err)
	opts := builder.Options{
		GeneratorOptions: inject.GeneratorOptions{ModulePath: modulePath, ModuleDir: moduleDir, OutputDir: out},
		CmdDir:           cmd,
	}

	files, err := inject.NewGenerator(g, opts.GeneratorOptions).Render(ctx)
	require.NoError(t, err)
	doors := builder.NewGenerator(spec, opts)
	rendered, err := doors.Render(ctx)
	require.NoError(t, err)
	for name, src := range rendered {
		files[name] = src
	}
	mains, err := doors.RenderCmd(ctx)
	require.NoError(t, err)

	const hint = "run go generate ./builder/examples/simple"
	onDisk := readGo(t, out, func(string) bool { return true })
	for name, src := range files {
		assert.Equal(t, string(src), onDisk[name], "gen/%s is stale: %s", name, hint)
		delete(onDisk, name)
	}
	assert.Empty(t, onDisk, "files in gen/ not produced by the generator: %s", hint)

	onDisk = readGo(t, cmd, func(name string) bool { return filepath.Base(name) == "main.go" })
	for name, src := range mains {
		assert.Equal(t, string(src), onDisk[name], "cmd/%s is stale: %s", name, hint)
		delete(onDisk, name)
	}
	assert.Empty(t, onDisk, "mains in cmd/ not produced by the generator: %s", hint)
}

// readGo returns the content of the files under dir that keep accepts, keyed
// by their slash-separated path relative to dir.
func readGo(t *testing.T, dir string, keep func(name string) bool) map[string]string {
	t.Helper()
	files := map[string]string{}
	require.NoError(t, filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !keep(p) {
			return err
		}
		rel, err := filepath.Rel(dir, p)
		if err != nil {
			return err
		}
		content, err := os.ReadFile(p)
		files[filepath.ToSlash(rel)] = string(content)
		return err
	}))
	return files
}
