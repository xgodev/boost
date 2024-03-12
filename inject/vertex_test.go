package inject

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type VertexTestSuite struct {
	suite.Suite
	graph *Graph[interface{}]
}

func TestVertexTestSuite(t *testing.T) {
	suite.Run(t, new(VertexTestSuite))
}

func (suite *VertexTestSuite) SetupSuite() {
	suite.graph = NewGraph[interface{}]()
	suite.graph.AddVertex("vertex1", "value1")
	suite.graph.AddVertex("vertex2", "value2")
	suite.graph.AddVertex("vertex3", "value3")
	suite.graph.AddEdge("vertex1", "vertex2")
	suite.graph.AddEdge("vertex2", "vertex3")
}

func (suite *VertexTestSuite) TestAdjacent() {
	testCases := []struct {
		name         string
		vertexKey    string
		expectedKeys []string
	}{
		{"Vertex With Adjacent", "vertex1", []string{"vertex2"}},
		{"Vertex Without Adjacent", "vertex3", []string{}},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			vertex := suite.graph.vertices[tc.vertexKey]
			adjacentVertices := vertex.Adjacent()

			var actualKeys []string
			for _, v := range adjacentVertices {
				actualKeys = append(actualKeys, v.Key)
			}

			suite.ElementsMatch(tc.expectedKeys, actualKeys)
		})
	}
}

func (suite *VertexTestSuite) TestIncoming() {
	testCases := []struct {
		name         string
		vertexKey    string
		expectedKeys []string
	}{
		{"Vertex With Incoming Edges", "vertex2", []string{"vertex1"}},
		{"Vertex Without Incoming Edges", "vertex1", []string{}},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			vertex := suite.graph.vertices[tc.vertexKey]
			incomingVertices := vertex.Incoming()

			var actualKeys []string
			for _, v := range incomingVertices {
				actualKeys = append(actualKeys, v.Key)
			}

			suite.ElementsMatch(tc.expectedKeys, actualKeys)
		})
	}
}
