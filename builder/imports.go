package builder

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"sort"

	"github.com/xgodev/boost/annotation"
	"github.com/xgodev/boost/inject"
	"github.com/xgodev/boost/model/errors"
	"golang.org/x/tools/go/ast/astutil"
)

// imports collects the imports of one generated file, one alias per path.
type imports struct {
	byPath map[string]string
}

func newImports() *imports {
	return &imports{byPath: map[string]string{}}
}

// add imports path under alias, or returns the alias it already has.
func (s *imports) add(alias, path string) string {
	if existing, ok := s.byPath[path]; ok {
		return existing
	}
	s.byPath[path] = alias
	return alias
}

// source imports a package the application wrote, under the alias inject
// gives source packages.
func (s *imports) source(name, path string) string {
	return s.add(inject.Alias("src", name, path), path)
}

// sorted returns the imports sorted by path.
func (s *imports) sorted() []inject.ImportData {
	list := make([]inject.ImportData, 0, len(s.byPath))
	for path, alias := range s.byPath {
		list = append(list, inject.ImportData{Alias: alias, Path: path})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Path < list[j].Path })
	return list
}

// qualify rewrites a type or function expression written in the file of
// entry so it compiles in a generated file: a name of the entry's package
// becomes alias.Name and pkg.Name uses the alias of the path the file imports
// as pkg. Struct field, parameter and method names stay as written.
func (s *imports) qualify(entry annotation.Entry, expr string) (string, error) {
	parsed, err := parser.ParseExpr(expr)
	if err != nil {
		return "", errors.NotValidf("%s: cannot parse %q: %v", declName(entry), expr, err)
	}

	var failure error
	rewritten := astutil.Apply(parsed, func(c *astutil.Cursor) bool {
		if failure != nil {
			return false
		}
		switch n := c.Node().(type) {
		case *ast.SelectorExpr:
			pkg, ok := n.X.(*ast.Ident)
			path, imported := "", false
			if ok {
				path, imported = entry.Imports[pkg.Name]
			}
			if !imported {
				failure = errors.NotValidf("%s: %s in %q is not a package the file imports", declName(entry), types.ExprString(n), expr)
				return false
			}
			if !token.IsExported(n.Sel.Name) {
				failure = errors.NotValidf("%s: %s in %q is not exported", declName(entry), types.ExprString(n), expr)
				return false
			}
			c.Replace(&ast.SelectorExpr{X: ast.NewIdent(s.source(pkg.Name, path)), Sel: n.Sel})
			return false
		case *ast.Ident:
			if _, isField := c.Parent().(*ast.Field); isField && c.Name() == "Names" {
				return false
			}
			if types.Universe.Lookup(n.Name) != nil {
				return false
			}
			if !token.IsExported(n.Name) {
				failure = errors.NotValidf("%s: %s in %q must be exported to be used by the generated code", declName(entry), n.Name, expr)
				return false
			}
			c.Replace(&ast.SelectorExpr{X: ast.NewIdent(s.source(entry.Package, entry.Path)), Sel: n})
			return false
		}
		return true
	}, nil)
	if failure != nil {
		return "", failure
	}

	return types.ExprString(rewritten.(ast.Expr)), nil
}
