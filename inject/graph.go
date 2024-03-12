package inject

import (
	"context"
	"fmt"
	"github.com/xgodev/boost/annotation"
	"github.com/xgodev/boost/errors"
	ustrings "github.com/xgodev/boost/utils/strings"
	"os"
	"strings"
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

type Component struct {
	Entry annotation.Entry
	An    Annotation
}

func NewGraphFromEntries(ctx context.Context, entries []annotation.Entry) (*Graph[Component], error) {

	out := make(map[string]Component)
	in := make(map[string][]Component)

	for _, entry := range entries {
		if !entry.IsFunc() {
			continue
		}

		if !isValidCombinedAnnotations(entry.Annotations) {
			log.Warnf("the combined annotations is invalid")
			continue
		}

		for _, ann := range entry.Annotations {
			if !isValidAnnotation(ann.Name) {
				log.Warnf("the annotation %s is invalid", ann.Name)
				continue
			}

			annType, _ := ParseAnnotationType(strings.ToUpper(ann.Name))

			a := Annotation{}

			err := ann.Decode(&a)
			if err != nil {
				return nil, err
			}

			switch annType {
			case AnnotationTypePROVIDE:

				if a.Index == nil {
					return nil, errors.NotValidf("the index parameter is required on the annotation %s in the entry %s.%s", ann.Name, entry.Path, entry.Func.Name)
				}

				index := *a.Index

				for i, res := range entry.Func.Results {
					if index != i {
						continue
					}

					id := xid(entry.Package, res.Type, a)
					if _, ok := out[id]; !ok {
						out[id] = Component{
							Entry: entry,
							An:    a,
						}
					}
				}

			case AnnotationTypeINJECT:

				if a.Index == nil {
					return nil, errors.NotValidf("the index parameter is required on the annotation %s in the entry %s.%s", ann.Name, entry.Path, entry.Func.Name)
				}

				index := *a.Index

				for i, param := range entry.Func.Parameters {
					if index != i {
						continue
					}

					id := xid(entry.Package, param.Type, a)

					if _, ok := out[id]; !ok {
						in[id] = make([]Component, 0)
					}
					in[id] = append(in[id], Component{
						Entry: entry,
						An:    a,
					})
				}

			case AnnotationTypeINVOKE:
			case AnnotationTypeMODULE:
			}
		}

	}

	graph := NewGraph[Component]()
	for id, ae := range out {
		graph.AddVertex(id, ae)
	}

	for id, aes := range in {

		if outAnnoEntry, ok := out[id]; ok {
			for _, inb := range aes {

				graph.AddVertex(gid(inb.Entry), inb)
				graph.AddEdge(gid(outAnnoEntry.Entry), gid(inb.Entry))

			}
		} else {
			return nil, errors.NotFoundf("provider not found for %s", id)
		}

	}

	return graph, nil
}

func gid(entry annotation.Entry) string {
	return strings.Join([]string{entry.Path, entry.Func.Name}, "_")
}

func xid(pkg string, tp string, ann Annotation) string {

	ref := false
	if strings.Contains(tp, "*") {
		ref = true
	}

	tp = strings.ReplaceAll(tp, "*", "")
	if !strings.Contains(tp, ".") {
		tp = strings.Join([]string{pkg, ".", tp}, "")
	}

	if ref {
		tp = strings.Join([]string{"*", tp}, "")
	}

	return strings.Join([]string{tp, ann.ID()}, "_")
}

func isValidCombinedAnnotations(annons []annotation.Annotation) bool {
	var all []string
	for _, ann := range annons {
		all = append(all, strings.ToUpper(ann.Name))
	}

	return !ustrings.SliceContainsAll(all, []string{AnnotationTypePROVIDE.String(), AnnotationTypeINJECT.String()})
}

func isValidAnnotation(value string) bool {
	if ustrings.SliceContains([]string{
		AnnotationTypeMODULE.String(),
		AnnotationTypePROVIDE.String(),
		AnnotationTypeINJECT.String(),
		AnnotationTypeINVOKE.String()},
		strings.ToUpper(value)) {
		return true
	}
	return false
}
