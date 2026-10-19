package builder

import (
	"go/ast"
	"go/parser"
	"strings"

	"github.com/xgodev/boost/annotation"
	"github.com/xgodev/boost/extra/graph"
	"github.com/xgodev/boost/inject"
	"github.com/xgodev/boost/model/errors"
)

// pluginPaths are the packages whose Plugin type a server of each kind is
// built with.
var pluginPaths = map[Kind]string{KindREST: boostEchoPath, KindGRPC: grpcServerPath, KindMCP: mcpServerPath}

// serverPlugins returns the inject constructors providing the []Plugin of
// srv and the name they provide it under: the value named after the server,
// or, for the rest, grpc and mcp servers, the default (unnamed) value. A
// default server given both is AlreadyExists: one of them would be dropped.
func serverPlugins(g *graph.Graph[inject.Component], srv Server) ([]inject.Component, string, error) {
	named := pluginProviders(g, srv, srv.Name)
	var unnamed []inject.Component
	if isDefault(srv) {
		unnamed = pluginProviders(g, srv, "")
	}

	switch {
	case len(named) > 0 && len(unnamed) > 0:
		return nil, "", errors.AlreadyExistsf("server %s of app %s takes the []Plugin named %s (%s) and the unnamed one (%s); provide one of them",
			srv.Name, srv.App, srv.Name, declNames(named), declNames(unnamed))
	case len(named) > 0:
		return named, srv.Name, nil
	}
	return unnamed, "", nil
}

// pluginProviders returns the inject constructors providing the ungrouped
// []Plugin of the kind of srv under name.
func pluginProviders(g *graph.Graph[inject.Component], srv Server, name string) []inject.Component {
	var plugins []inject.Component
	for _, key := range g.SortedKeys() {
		comp := g.Vertices[key].Value
		for idx, ann := range comp.Provides {
			if ann.Group == "" && ann.Name == name && isPluginSlice(comp.Entry, comp.Result(idx), pluginPaths[srv.Kind]) {
				plugins = append(plugins, comp)
				break
			}
		}
	}
	return plugins
}

func declNames(comps []inject.Component) string {
	names := make([]string, 0, len(comps))
	for _, comp := range comps {
		names = append(names, declName(comp.Entry))
	}
	return strings.Join(names, ", ")
}

// isPluginSlice reports whether tp, written in the file of entry, is
// []pkg.Plugin with pkg imported from pkgPath.
func isPluginSlice(entry annotation.Entry, tp, pkgPath string) bool {
	expr, err := parser.ParseExpr(tp)
	if err != nil {
		return false
	}
	slice, ok := expr.(*ast.ArrayType)
	if !ok || slice.Len != nil {
		return false
	}
	sel, ok := slice.Elt.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Plugin" {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	return ok && entry.Imports[pkg.Name] == pkgPath
}
