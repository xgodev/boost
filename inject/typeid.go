package inject

import (
	"go/ast"
	"go/parser"
	"go/types"

	"golang.org/x/tools/go/ast/astutil"
)

// qualifyType rewrites the type expression tp, written inside package pkg,
// so every package-local type name carries the package qualifier. "*Loren"
// declared in package simple and "*simple.Loren" seen from another package
// end up with the same id. Predeclared types (string, error, ...) are kept.
func qualifyType(pkg, tp string) string {
	expr, err := parser.ParseExpr(tp)
	if err != nil {
		return tp
	}

	out := astutil.Apply(expr, func(c *astutil.Cursor) bool {
		switch n := c.Node().(type) {
		case *ast.SelectorExpr:
			// already qualified (pkg.Type)
			return false
		case *ast.Ident:
			// parameter, result and struct field names are not types
			if c.Name() == "Names" || isPredeclaredType(n.Name) {
				return false
			}
			c.Replace(&ast.SelectorExpr{X: ast.NewIdent(pkg), Sel: ast.NewIdent(n.Name)})
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
