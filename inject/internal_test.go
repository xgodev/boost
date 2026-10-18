package inject

import (
	"context"
	"go/parser"
	"go/token"
	"path"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost/annotation"
)

// importAllowed reports whether Go lets package importer import imported: a
// path with an internal element is importable only from the tree rooted at
// the parent of its last internal element.
func importAllowed(importer, imported string) bool {
	elems := strings.Split(imported, "/")
	for i := len(elems) - 1; i >= 0; i-- {
		if elems[i] == "internal" {
			parent := strings.Join(elems[:i], "/")
			return importer == parent || strings.HasPrefix(importer, parent+"/")
		}
	}
	return true
}

// assertImportsAllowed fails when a file of files, rendered under the import
// path base, imports a package that Go's internal rule hides from it.
func assertImportsAllowed(t *testing.T, base string, files map[string][]byte) {
	t.Helper()
	for name, src := range files {
		f, err := parser.ParseFile(token.NewFileSet(), name, src, parser.ImportsOnly)
		require.NoError(t, err, "%s is not valid Go", name)
		importer := path.Join(base, path.Dir(name))
		for _, imp := range f.Imports {
			p, err := strconv.Unquote(imp.Path.Value)
			require.NoError(t, err)
			assert.True(t, importAllowed(importer, p), "%s (package %s) cannot import %s", name, importer, p)
		}
	}
}

func TestImportAllowed(t *testing.T) {
	assert.True(t, importAllowed("m/gen", "m/gen/a"))
	assert.True(t, importAllowed("m/x", "m/internal/a"))
	assert.True(t, importAllowed("m", "m/internal/a"))
	assert.False(t, importAllowed("m/gen", "m/gen/m/internal/a"))
	assert.False(t, importAllowed("m/x", "m/x/internal/y/internal/z"))
	assert.True(t, importAllowed("m/x/internal/y/w", "m/x/internal/y/internal/z"))
}

// TestGenerator_RenderInternalPackages renders the components of the xgodev
// layout, where every package lives under internal/, and checks that each
// generated file may import what it imports.
func TestGenerator_RenderInternalPackages(t *testing.T) {
	stock := map[string]string{"stock": testModule + "/internal/pkg/domain/repository/stock"}
	memory := at("internal/pkg/datastore/stock/memory", "New", nil, []string{"stock.Repository"}, provide("index=0"), deflt())
	memory.Imports = stock
	mongodb := at("internal/pkg/datastore/stock/mongodb", "New", nil, []string{"stock.Repository"}, provide("index=0"))
	mongodb.Imports = stock
	clock := at("internal/pkg/clock", "New", nil, []string{"*Clock"}, provide("index=0"))
	service := at("internal/pkg/domain/service/stock/reserver/local", "New", []string{"stock.Repository", "*clock.Clock"}, []string{"*Service"}, provide("index=0"))
	service.Imports = map[string]string{"stock": stock["stock"], "clock": testModule + "/internal/pkg/clock"}

	files, err := NewGenerator(newTestGraph(t, []annotation.Entry{memory, mongodb, clock, service}), testOptions("/mod")).Render(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, files)

	assertImportsAllowed(t, testModule+"/gen", files)
}
