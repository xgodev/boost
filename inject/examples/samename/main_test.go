package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/inject"
	"github.com/xgodev/boost/inject/examples/samename/app"
	"github.com/xgodev/boost/inject/examples/samename/gen"
	"go.uber.org/fx"
)

// Each boost package gets its own Options: the two types share a package
// name but not an import path (#202).
func TestGeneratedModuleWiresTheApp(t *testing.T) {
	boost.Start() // loads the config keys the generated modules registered

	var subjects *app.Subjects
	require.NoError(t, fx.New(gen.Module(), fx.Populate(&subjects), fx.NopLogger).Err())

	assert.Equal(t, &app.Subjects{Received: "precos-entrada", Changed: "precos-alteracao"}, subjects)
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

	const hint = "gen/ is stale: run go generate ./inject/examples/samename"
	for name, src := range files {
		assert.Equal(t, string(src), onDisk[name], "%s: %s", name, hint)
		delete(onDisk, name)
	}
	assert.Empty(t, onDisk, "files not produced by the generator: %s", hint)
}
