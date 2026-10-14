package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost/inject"
	"github.com/xgodev/boost/inject/examples/simple/app"
	"github.com/xgodev/boost/inject/examples/simple/gen"
	"go.uber.org/fx"
)

func TestGeneratedModuleWiresTheApp(t *testing.T) {
	var board *app.Board
	require.NoError(t, fx.New(gen.Module(), fx.Populate(&board), fx.NopLogger).Err())

	assert.Equal(t, []string{
		"Good morning, Alice!",
		"Hello, Alice!",
		"Good morning, Bob!",
		"Hello, Bob!",
	}, board.Lines)
}

// TestGeneratedCodeIsUpToDate fails when gen/ differs from what the
// generator produces for the annotations in this example. Run go generate.
func TestGeneratedCodeIsUpToDate(t *testing.T) {
	ctx := context.Background()

	modulePath, moduleDir, err := inject.FindModule(".")
	require.NoError(t, err)

	entries, err := inject.CollectEntries(".")
	require.NoError(t, err)

	g, err := inject.NewGraphFromEntries(ctx, entries)
	require.NoError(t, err)

	out, err := filepath.Abs("gen")
	require.NoError(t, err)

	files, err := inject.NewGenerator(g, inject.GeneratorOptions{
		ModulePath: modulePath,
		ModuleDir:  moduleDir,
		OutputDir:  out,
	}).Render(ctx)
	require.NoError(t, err)

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

	const hint = "gen/ is stale: run go generate ./inject/examples/simple"
	for name, src := range files {
		assert.Equal(t, string(src), onDisk[name], "%s: %s", name, hint)
		delete(onDisk, name)
	}
	assert.Empty(t, onDisk, "files not produced by the generator: %s", hint)
}
