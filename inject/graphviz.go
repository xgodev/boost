package inject

import (
	"fmt"
	"github.com/xgodev/boost/extra/graph"
	"os"
)

// ExportInjectGraphToGraphviz writes the dependency injection graph in Graphviz's DOT format.
// Each module in the graph is represented with a specific color and shape, with edges illustrating
// dependencies between modules, configured according to the annotation type (e.g., @Provide, @Inject).
//
// Parameters:
// - g: The dependency graph containing injection components.
// - filename: The output file name to store the graph in DOT format.
//
// Returns:
// - An error if any failure occurs while writing to the file.
func ExportInjectGraphToGraphviz(g *graph.Graph[Component], filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating DOT file: %w", err)
	}
	defer file.Close()

	// Basic DOT file header with global styling configurations
	_, err = file.WriteString(`digraph InjectGraph {
        graph [splines=true, overlap=false];
        node [fontname="Arial", shape=box, style=filled];
        edge [color="#606060"];
    `)
	if err != nil {
		return fmt.Errorf("error writing header to DOT file: %w", err)
	}

	// Iterate over vertices to define nodes in the graph
	for _, vertex := range g.Vertices {
		component := vertex.Value
		nodeID := fmt.Sprintf("%s_%s", component.Entry.Package, component.Entry.Func.Name)

		// Define node color and shape based on the presence of `Provide` and `Inject` annotations
		nodeColor := "#A0A0FF" // Default color
		nodeShape := "box"     // Default shape for modules

		// Adjust style for @Provide or @Inject based on component annotations
		for _, ann := range component.Entry.Annotations {
			if ann.Name == "Provide" {
				nodeColor = "#90EE90" // Green for providers
				nodeShape = "ellipse"
				break
			} else if ann.Name == "Inject" {
				nodeColor = "#FFD700" // Yellow for injectors
				nodeShape = "octagon"
				break
			}
		}

		// Write the node to the DOT file with customized attributes
		_, err = file.WriteString(fmt.Sprintf(
			"\"%s\" [label=\"%s\", shape=%s, fillcolor=\"%s\"];\n",
			nodeID, component.Entry.Func.Name, nodeShape, nodeColor,
		))
		if err != nil {
			return fmt.Errorf("error writing node to DOT file: %w", err)
		}

		// Iterate over dependencies and write edges between nodes
		for _, edge := range g.Edges[vertex.Key] {
			dependencyID := fmt.Sprintf("%s_%s", edge.Value.Entry.Package, edge.Value.Entry.Func.Name)
			_, err = file.WriteString(fmt.Sprintf("\"%s\" -> \"%s\";\n", nodeID, dependencyID))
			if err != nil {
				return fmt.Errorf("error writing edge to DOT file: %w", err)
			}
		}
	}

	// Finalize the DOT file
	_, err = file.WriteString("}\n")
	if err != nil {
		return fmt.Errorf("error finalizing DOT file: %w", err)
	}

	return nil
}
