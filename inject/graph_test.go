package inject

import (
	"context"
	"github.com/xgodev/boost/annotation"
	"github.com/xgodev/boost/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
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
			suite.Contains(g.vertices, tc.key)
			suite.Equal(tc.value, g.vertices[tc.key].Value)
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

			suite.Len(g.edges[tc.from], 1)
			suite.Equal(tc.to, g.edges[tc.from][0].Key)
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

type NewGraphFromEntriesTestSuite struct {
	suite.Suite
	testData map[string][]annotation.Entry
}

func TestNewGraphFromEntriesTestSuite(t *testing.T) {
	suite.Run(t, new(GraphTestSuite))
}

func (suite *NewGraphFromEntriesTestSuite) SetupSuite() {

	suite.testData = make(map[string][]annotation.Entry)

	files, err := filepath.Glob("testdata/inject/mkgraph/*.yaml")
	suite.Require().NoError(err)

	for _, file := range files {
		id := filepath.Base(file)
		entries, err := suite.loadEntriesFromYAML(file)
		suite.Require().NoError(err)

		suite.testData[id] = entries
	}
}

func (suite *NewGraphFromEntriesTestSuite) loadEntriesFromYAML(filename string) ([]annotation.Entry, error) {
	var entries []annotation.Entry

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &entries)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func (suite *NewGraphFromEntriesTestSuite) TestNewGraphFromEntries() {
	testCases := []struct {
		name      string
		id        string // ID do cenário de teste
		expectErr bool
		errType   error
	}{
		{
			name:      "valid",
			id:        "1_valid.yaml",
			expectErr: false,
		},
		{
			name:      "provider not found",
			id:        "2_provider_notfound.yaml",
			expectErr: true,
			errType:   errors.NotFoundf("tipo de erro esperado"),
		},
		{
			name:      "provide index not found",
			id:        "3_provide_index_notfound.yaml",
			expectErr: true,
			errType:   errors.NotValidf("tipo de erro esperado"),
		},
		{
			name:      "inject index not found",
			id:        "4_inject_index_notfound.yaml",
			expectErr: true,
			errType:   errors.NotValidf("tipo de erro esperado"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			entries := suite.testData[tc.id]
			_, err := NewGraphFromEntries(context.Background(), entries)

			if tc.expectErr {
				suite.Error(err)
				suite.IsType(tc.errType, err)
			} else {
				suite.NoError(err)
			}
		})
	}
}
