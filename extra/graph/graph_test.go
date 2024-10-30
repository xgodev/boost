package graph

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type GraphTestSuite struct {
	suite.Suite
}

func TestGraphTestSuite(t *testing.T) {
	suite.Run(t, new(GraphTestSuite))
}

func (suite *GraphTestSuite) TestAddVertex() {
	testCases := []struct {
		name  string
		key   string
		value string
	}{
		{"Add Single Vertex", "vertex1", "value1"},
		{"Add Another Vertex", "vertex2", "value2"},
	}

	g := NewGraph[string]()
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			g.AddVertex(tc.key, tc.value)
			suite.Contains(g.Vertices, tc.key)
			suite.Equal(tc.value, g.Vertices[tc.key].Value)
		})
	}
}

func (suite *GraphTestSuite) TestAddEdge() {
	testCases := []struct {
		name       string
		from, to   string
		setupGraph func(g *Graph[string])
	}{
		{
			name: "Add Edge Between Existing Vertices",
			from: "vertex1",
			to:   "vertex2",
			setupGraph: func(g *Graph[string]) {
				g.AddVertex("vertex1", "value1")
				g.AddVertex("vertex2", "value2")
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			g := NewGraph[string]()
			tc.setupGraph(g)
			g.AddEdge(tc.from, tc.to)

			suite.Len(g.Edges[tc.from], 1)
			suite.Equal(tc.to, g.Edges[tc.from][0].Key)
		})
	}
}

func (suite *GraphTestSuite) TestVerticesWithNoIncomingEdges() {
	testCases := []struct {
		name       string
		setupGraph func(g *Graph[string])
		expected   []string
	}{
		{
			name: "Single Vertex With No Incoming Edge",
			setupGraph: func(g *Graph[string]) {
				g.AddVertex("vertex1", "value1")
			},
			expected: []string{"vertex1"},
		},
		{
			name: "One Vertex With Incoming Edge",
			setupGraph: func(g *Graph[string]) {
				g.AddVertex("vertex1", "value1")
				g.AddVertex("vertex2", "value2")
				g.AddEdge("vertex1", "vertex2")
			},
			expected: []string{"vertex1"},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			g := NewGraph[string]()
			tc.setupGraph(g)

			vertices := g.VerticesWithNoIncomingEdges()
			suite.Len(vertices, len(tc.expected))
			for _, expectedKey := range tc.expected {
				found := false
				for _, vertex := range vertices {
					if vertex.Key == expectedKey {
						found = true
						break
					}
				}
				suite.True(found)
			}
		})
	}
}
