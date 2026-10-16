package inject

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/xgodev/boost/extra/graph"
	"github.com/xgodev/boost/model/errors"
)

// ExportInjectGraphToGraphviz writes the dependency graph to filename in
// Graphviz DOT format. Providers are green ellipses, invokers yellow
// octagons, and an edge goes from a provider to each component using it.
func ExportInjectGraphToGraphviz(g *graph.Graph[Component], filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, errors.Internalf("error creating DOT file %s", filename))
	}
	defer file.Close()

	return WriteGraphviz(file, g)
}

// WriteGraphviz writes the dependency graph to w in Graphviz DOT format,
// in a stable order.
func WriteGraphviz(w io.Writer, g *graph.Graph[Component]) error {
	if _, err := fmt.Fprint(w, `digraph InjectGraph {
	graph [splines=true, overlap=false];
	node [fontname="Arial", style=filled];
	edge [color="#606060"];
`); err != nil {
		return errors.Wrap(err, errors.Internalf("error writing DOT header"))
	}

	for _, key := range g.SortedKeys() {
		comp := g.Vertices[key].Value

		shape, color := "ellipse", "#90EE90"
		if comp.Type == AnnotationTypeINVOKE {
			shape, color = "octagon", "#FFD700"
		}

		if _, err := fmt.Fprintf(w, "\t%q [label=%q, shape=%s, fillcolor=%q];\n",
			key, comp.Entry.Package+"."+comp.Entry.Func.Name, shape, color); err != nil {
			return errors.Wrap(err, errors.Internalf("error writing DOT node"))
		}

		var deps []string
		for _, adj := range g.Edges[key] {
			deps = append(deps, adj.Key)
		}
		sort.Strings(deps)
		for _, dep := range deps {
			if _, err := fmt.Fprintf(w, "\t%q -> %q;\n", key, dep); err != nil {
				return errors.Wrap(err, errors.Internalf("error writing DOT edge"))
			}
		}
	}

	if _, err := fmt.Fprint(w, "}\n"); err != nil {
		return errors.Wrap(err, errors.Internalf("error writing DOT footer"))
	}
	return nil
}
