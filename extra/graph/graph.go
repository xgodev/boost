package graph

import (
	"fmt"
	"github.com/xgodev/boost/wrapper/log"
	"os"
)

// Graph represents a generic graph structure with vertices of any type.
// It includes maps for vertices, incoming edges, and edges for efficient graph operations.
type Graph[T any] struct {
	vertices      map[string]*Vertex[T]   // Map of vertices in the graph.
	incomingEdges map[string]int          // Map of incoming edge counts per vertex.
	edges         map[string][]*Vertex[T] // Map of edges represented as adjacency lists.
}

// NewGraph creates and returns a new instance of Graph.
func NewGraph[T any]() *Graph[T] {
	return &Graph[T]{
		vertices:      make(map[string]*Vertex[T]),
		incomingEdges: make(map[string]int),
		edges:         make(map[string][]*Vertex[T]),
	}
}

// AddVertex adds a new vertex with the specified key and value to the graph.
// If the vertex already exists, it logs a warning and does not overwrite it.
func (g *Graph[T]) AddVertex(key string, value T) {
	if _, exists := g.vertices[key]; exists {
		log.Warnf("vertex %s already exists", key)
		return
	}

	vertex := &Vertex[T]{Key: key, Value: value, graph: g}
	g.vertices[key] = vertex
	g.incomingEdges[key] = 0

	log.Debugf("vertex added: %s", key)
}

// AddEdge adds a directed edge from one vertex to another.
// If either vertex does not exist, it logs a warning and does not add the edge.
func (g *Graph[T]) AddEdge(fromKey, toKey string) {
	_, fromExists := g.vertices[fromKey]
	toVertex, toExists := g.vertices[toKey]

	if !fromExists {
		log.Warnf("from vertex not found. %v ", fromKey)
		return
	}

	if !toExists {
		log.Warnf("to vertex not found. %v ", toKey)
		return
	}

	if fromKey == toKey {
		log.Warnf("warning: cannot add edge to the same vertex: %v", fromKey)
		return
	}

	for _, v := range g.edges[fromKey] {
		if v.Key == toKey {
			return
		}
	}

	g.edges[fromKey] = append(g.edges[fromKey], toVertex)
	g.incomingEdges[toKey]++
	log.Debugf("edge added from %v to %v", fromKey, toKey)
}

// VerticesWithNoIncomingEdges returns a list of vertices with no incoming edges.
func (g *Graph[T]) VerticesWithNoIncomingEdges() []*Vertex[T] {
	var vertices []*Vertex[T]
	for key, count := range g.incomingEdges {
		if count == 0 {
			vertices = append(vertices, g.vertices[key])
		}
	}
	return vertices
}

func (g *Graph[T]) Print() {
	for _, vertex := range g.vertices {
		log.Infof("%v (%v) -> ", vertex.Key, vertex.Value)
		for _, edge := range g.edges[vertex.Key] {
			log.Infof("%v ", edge.Value)
		}
	}
}

func (g *Graph[T]) ExportToGraphviz(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString("digraph G {\n")
	if err != nil {
		return err
	}

	for key, _ := range g.vertices {
		_, err = file.WriteString(fmt.Sprintf("\t\"%s\" [label=\"%s\"];\n", key, key))
		if err != nil {
			return err
		}

		for _, edge := range g.edges[key] {
			_, err = file.WriteString(fmt.Sprintf("\t\"%s\" -> \"%s\";\n", key, edge.Key))
			if err != nil {
				return err
			}
		}
	}

	_, err = file.WriteString("}\n")
	if err != nil {
		return err
	}

	return nil
}
