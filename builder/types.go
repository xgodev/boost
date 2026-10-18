package builder

import (
	"go/ast"
	"go/parser"
	"go/types"
	"strings"

	"github.com/xgodev/boost/annotation"
)

const (
	echoPath        = "github.com/labstack/echo/v4"
	cloudeventsPath = "github.com/cloudevents/sdk-go/v2"
	eventPath       = "github.com/cloudevents/sdk-go/v2/event"
	mcpPath         = "github.com/modelcontextprotocol/go-sdk/mcp"
)

// typeRef is a named type or function: its package import path and name.
type typeRef struct {
	Path string
	Name string
}

// resolve tells which package declares the type or function written as Name
// or pkg.Name in the file of entry. ok is false for any other expression,
// predeclared identifiers and packages the file does not import.
func resolve(entry annotation.Entry, expr string) (ref typeRef, ok bool) {
	parsed, err := parser.ParseExpr(expr)
	if err != nil {
		return typeRef{}, false
	}
	switch e := parsed.(type) {
	case *ast.Ident:
		if types.Universe.Lookup(e.Name) != nil {
			return typeRef{}, false
		}
		return typeRef{Path: entry.Path, Name: e.Name}, true
	case *ast.SelectorExpr:
		pkg, isIdent := e.X.(*ast.Ident)
		if !isIdent {
			return typeRef{}, false
		}
		importPath, imported := entry.Imports[pkg.Name]
		if !imported {
			return typeRef{}, false
		}
		return typeRef{Path: importPath, Name: e.Sel.Name}, true
	}
	return typeRef{}, false
}

// isPointer reports whether expr, written in the file of entry, is *ref.
func isPointer(entry annotation.Entry, expr string, ref typeRef) bool {
	elem, ok := strings.CutPrefix(strings.TrimSpace(expr), "*")
	return ok && isType(entry, elem, ref)
}

// isType reports whether expr, written in the file of entry, is one of refs.
func isType(entry annotation.Entry, expr string, refs ...typeRef) bool {
	got, ok := resolve(entry, expr)
	if !ok {
		return false
	}
	for _, ref := range refs {
		if got == ref {
			return true
		}
	}
	return false
}
