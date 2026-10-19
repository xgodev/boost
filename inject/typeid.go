package inject

import (
	"go/ast"
	"go/parser"
	"go/types"

	"github.com/xgodev/boost/annotation"
	"golang.org/x/tools/go/ast/astutil"
)

// qualifyType rewrites the type expression tp, written in the file of entry,
// so every type name carries the import path of its package. "*Loren"
// declared in package example.com/app/simple and "*simple.Loren" seen from a
// file importing that package end up with the same id, and two packages that
// share a name are told apart. A qualifier the file does not import is kept as
// written. Predeclared types (string, error, ...) are kept.
func qualifyType(entry annotation.Entry, tp string) string {
	return rewriteType(tp, func(name string) string {
		if path, ok := entry.Imports[name]; ok {
			return path
		}
		return name
	}, entry.Path)
}

// qualifyName is qualifyType by package name: the id Provider matches on.
func qualifyName(pkg, tp string) string {
	return rewriteType(tp, func(name string) string { return name }, pkg)
}

// rewriteType qualifies every type name of tp: a qualified one (pkg.Type) by
// what qualifier returns for pkg, a local one by local.
func rewriteType(tp string, qualifier func(string) string, local string) string {
	expr, err := parser.ParseExpr(tp)
	if err != nil {
		return tp
	}

	out := astutil.Apply(expr, func(c *astutil.Cursor) bool {
		switch n := c.Node().(type) {
		case *ast.SelectorExpr:
			if pkg, ok := n.X.(*ast.Ident); ok {
				c.Replace(&ast.SelectorExpr{X: ast.NewIdent(qualifier(pkg.Name)), Sel: n.Sel})
			}
			return false
		case *ast.Ident:
			// parameter, result and struct field names are not types
			if c.Name() == "Names" || isPredeclaredType(n.Name) {
				return false
			}
			c.Replace(&ast.SelectorExpr{X: ast.NewIdent(local), Sel: ast.NewIdent(n.Name)})
			return false
		}
		return true
	}, nil)

	return types.ExprString(out.(ast.Expr))
}

// sliceElem returns the element type of a slice type expression, used to
// match a `[]T` group parameter with the providers of T.
func sliceElem(tp string) (string, bool) {
	expr, err := parser.ParseExpr(tp)
	if err != nil {
		return "", false
	}
	arr, ok := expr.(*ast.ArrayType)
	if !ok || arr.Len != nil {
		return "", false
	}
	return types.ExprString(arr.Elt), true
}

func isPredeclaredType(name string) bool {
	_, ok := types.Universe.Lookup(name).(*types.TypeName)
	return ok
}

// withImports returns entries with the imports a hand-built entry leaves out:
// a qualifier its function types use that is not in Imports resolves to the
// package of the entries with that name, when only one has it. The collector
// fills every import of a file, so its entries come back unchanged.
func withImports(entries []annotation.Entry) []annotation.Entry {
	paths := map[string]map[string]bool{} // package name -> import paths
	for _, entry := range entries {
		if paths[entry.Package] == nil {
			paths[entry.Package] = map[string]bool{}
		}
		paths[entry.Package][entry.Path] = true
	}

	out := make([]annotation.Entry, len(entries))
	for i, entry := range entries {
		copied := false
		for _, name := range qualifiers(entry) {
			if _, ok := entry.Imports[name]; ok || len(paths[name]) != 1 {
				continue
			}
			if !copied {
				// entries of a file share the map
				imports := make(map[string]string, len(entry.Imports)+1)
				for k, v := range entry.Imports {
					imports[k] = v
				}
				entry.Imports, copied = imports, true
			}
			for p := range paths[name] {
				entry.Imports[name] = p
			}
		}
		out[i] = entry
	}
	return out
}

// qualifiers returns the package qualifiers used by the parameter and result
// types of entry.
func qualifiers(entry annotation.Entry) []string {
	var names []string
	for _, list := range [][]annotation.EntryFuncType{entry.Func.Parameters, entry.Func.Results} {
		for _, param := range list {
			expr, err := parser.ParseExpr(param.Type)
			if err != nil {
				continue
			}
			ast.Inspect(expr, func(n ast.Node) bool {
				if sel, ok := n.(*ast.SelectorExpr); ok {
					if pkg, ok := sel.X.(*ast.Ident); ok {
						names = append(names, pkg.Name)
					}
					return false
				}
				return true
			})
		}
	}
	return names
}
