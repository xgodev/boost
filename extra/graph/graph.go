package graph

import (
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/xgodev/boost/wrapper/log"
)

// ErrCycle is returned by TopologicalSort when the graph has a cycle.
var ErrCycle = errors.New("graph has a cycle")

// Graph represents a generic graph structure with Vertices of any type.
// It includes maps for Vertices, incoming Edges, and Edges for efficient graph operations.
type Graph[T any] struct {
	Vertices      map[string]*Vertex[T]   // Map of Vertices in the graph.
	IncomingEdges map[string]int          // Map of incoming edge counts per vertex.
	Edges         map[string][]*Vertex[T] // Map of Edges represented as adjacency lists.
}

// NewGraph creates and returns a new instance of Graph.
func NewGraph[T any]() *Graph[T] {
	return &Graph[T]{
		Vertices:      make(map[string]*Vertex[T]),
		IncomingEdges: make(map[string]int),
		Edges:         make(map[string][]*Vertex[T]),
	}
}

// AddVertex adds a new vertex with the specified key and value to the graph.
// If the vertex already exists, it logs a warning and does not overwrite it.
func (g *Graph[T]) AddVertex(key string, value T) {
	if _, exists := g.Vertices[key]; exists {
		log.Warnf("vertex %s already exists", key)
		return
	}

	vertex := &Vertex[T]{Key: key, Value: value, graph: g}
	g.Vertices[key] = vertex
	g.IncomingEdges[key] = 0

	log.Debugf("vertex added: %s", key)
}

// AddEdge adds a directed edge from one vertex to another.
// If either vertex does not exist, it logs a warning and does not add the edge.
func (g *Graph[T]) AddEdge(fromKey, toKey string) {
	_, fromExists := g.Vertices[fromKey]
	toVertex, toExists := g.Vertices[toKey]

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

	for _, v := range g.Edges[fromKey] {
		if v.Key == toKey {
			return
		}
	}

	g.Edges[fromKey] = append(g.Edges[fromKey], toVertex)
	g.IncomingEdges[toKey]++
	log.Debugf("edge added from %v to %v", fromKey, toKey)
}

// VerticesWithNoIncomingEdges returns a list of Vertices with no incoming Edges.
func (g *Graph[T]) VerticesWithNoIncomingEdges() []*Vertex[T] {
	var vertices []*Vertex[T]
	for key, count := range g.IncomingEdges {
		if count == 0 {
			vertices = append(vertices, g.Vertices[key])
		}
	}
	return vertices
}

func (g *Graph[T]) Print() {
	for _, vertex := range g.Vertices {
		log.Infof("%v (%v) -> ", vertex.Key, vertex.Value)
		for _, edge := range g.Edges[vertex.Key] {
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

	// Cabeçalho básico do arquivo DOT com algumas configurações de estilo globais
	_, err = file.WriteString(`digraph G {
        graph [splines=true, overlap=false];
        node [shape=box, style=filled, fontname="Arial"];
        edge [color="#606060"];
    `)
	if err != nil {
		return err
	}

	// Itera sobre os nós (vértices) do grafo
	for key := range g.Vertices {
		// Configurações básicas para os nós, mantendo o grafo genérico
		nodeColor := "#A0A0FF" // Cor padrão
		nodeShape := "box"     // Forma padrão

		// Escreve os nós no arquivo, mantendo uma abordagem genérica
		_, err = file.WriteString(fmt.Sprintf(
			"\t\"%s\" [label=\"%s\", shape=%s, fillcolor=\"%s\"];\n",
			key, key, nodeShape, nodeColor))
		if err != nil {
			return err
		}

		// Escreve as arestas (conexões) entre os nós
		for _, edge := range g.Edges[key] {
			edgeStyle := "solid"   // Estilo padrão das arestas
			edgeColor := "#000000" // Cor padrão

			// Escreve as arestas entre os nós
			_, err = file.WriteString(fmt.Sprintf(
				"\t\"%s\" -> \"%s\" [style=%s, color=\"%s\"];\n",
				key, edge.Key, edgeStyle, edgeColor))
			if err != nil {
				return err
			}
		}
	}

	// Finaliza o arquivo DOT
	_, err = file.WriteString("}\n")
	if err != nil {
		return err
	}

	return nil
}

// SortedKeys returns the keys of all vertices in lexical order, so callers
// that iterate the graph produce the same output on every run.
func (g *Graph[T]) SortedKeys() []string {
	keys := make([]string, 0, len(g.Vertices))
	for key := range g.Vertices {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// TopologicalSort returns the vertices ordered so that every vertex comes
// after all vertices with an edge into it. Ties are broken by key, so the
// order is deterministic. When the graph has a cycle, it returns
// ErrCycle and the keys of the vertices that are part of (or depend on)
// the cycle.
func (g *Graph[T]) TopologicalSort() ([]*Vertex[T], []string, error) {
	incoming := make(map[string]int, len(g.IncomingEdges))
	for key, count := range g.IncomingEdges {
		incoming[key] = count
	}

	var ready []string
	for _, key := range g.SortedKeys() {
		if incoming[key] == 0 {
			ready = append(ready, key)
		}
	}

	sorted := make([]*Vertex[T], 0, len(g.Vertices))
	for len(ready) > 0 {
		key := ready[0]
		ready = ready[1:]
		sorted = append(sorted, g.Vertices[key])

		var next []string
		for _, adj := range g.Edges[key] {
			incoming[adj.Key]--
			if incoming[adj.Key] == 0 {
				next = append(next, adj.Key)
			}
		}
		sort.Strings(next)
		ready = append(ready, next...)
		sort.Strings(ready)
	}

	if len(sorted) == len(g.Vertices) {
		return sorted, nil, nil
	}

	var remaining []string
	for _, key := range g.SortedKeys() {
		if incoming[key] > 0 {
			remaining = append(remaining, key)
		}
	}
	return nil, remaining, ErrCycle
}
