package inject

// Vertex represents a node in a graph with a key and a value.
// It also holds a reference to the graph it belongs to for graph operations.
type Vertex[T any] struct {
	Key   string    // Unique identifier of the vertex.
	Value T         // Value stored in the vertex.
	graph *Graph[T] // Reference to the graph the vertex belongs to.
}

// Adjacent returns a list of adjacent vertices to this vertex.
// It returns nil if the vertex is not part of a graph.
func (v *Vertex[T]) Adjacent() []*Vertex[T] {
	if v.graph != nil {
		return v.graph.edges[v.Key]
	}
	return nil
}

// Incoming returns a list of vertices with edges incoming to this vertex.
// It iterates over all edges in the graph to find incoming connections.
func (v *Vertex[T]) Incoming() []*Vertex[T] {
	var incomingVertices []*Vertex[T]
	for key, edges := range v.graph.edges {
		for _, edge := range edges {
			if edge.Key == v.Key {
				incomingVertices = append(incomingVertices, v.graph.vertices[key])
			}
		}
	}
	return incomingVertices
}
