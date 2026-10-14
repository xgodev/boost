package inject

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/xgodev/boost/annotation"
	"github.com/xgodev/boost/extra/graph"
	"github.com/xgodev/boost/model/errors"
	"gopkg.in/yaml.v3"
)

type NewGraphFromEntriesTestSuite struct {
	suite.Suite
	testData map[string][]annotation.Entry
}

func TestNewGraphFromEntriesTestSuite(t *testing.T) {
	suite.Run(t, new(NewGraphFromEntriesTestSuite))
}

func (suite *NewGraphFromEntriesTestSuite) SetupSuite() {
	suite.testData = make(map[string][]annotation.Entry)

	files, err := filepath.Glob("testdata/graph/*.yaml")
	suite.Require().NoError(err)
	suite.Require().NotEmpty(files)

	for _, file := range files {
		data, err := os.ReadFile(file)
		suite.Require().NoError(err)

		var entries []annotation.Entry
		suite.Require().NoError(yaml.Unmarshal(data, &entries))
		suite.testData[filepath.Base(file)] = entries
	}
}

func (suite *NewGraphFromEntriesTestSuite) TestFixtures() {
	testCases := []struct {
		name  string
		file  string
		isErr func(error) bool
	}{
		{name: "valid", file: "1_success.yaml"},
		{name: "provider not found", file: "2_provider_notfound.yaml", isErr: errors.IsNotFound},
		{name: "provide index not found", file: "3_provide_index_notfound.yaml", isErr: errors.IsNotValid},
		{name: "inject index not found", file: "4_inject_index_notfound.yaml", isErr: errors.IsNotValid},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			entries, ok := suite.testData[tc.file]
			suite.Require().True(ok, "fixture %s not loaded", tc.file)

			_, err := NewGraphFromEntries(context.Background(), entries)
			if tc.isErr == nil {
				suite.NoError(err)
				return
			}
			suite.Error(err)
			suite.True(tc.isErr(err), "unexpected error type: %v", err)
		})
	}
}

func (suite *NewGraphFromEntriesTestSuite) TestValidFixtureEdges() {
	g, err := NewGraphFromEntries(context.Background(), suite.testData["1_success.yaml"])
	suite.Require().NoError(err)

	const (
		annotated = "github.com/jpfaria/tests/annotated."
		simple    = "github.com/xgodev/boost/inject/examples/simple."
	)

	suite.ElementsMatch([]string{simple + "Bar", simple + "Foz"}, edgesOf(g.Edges[simple+"FooBar"]))
	suite.ElementsMatch([]string{simple + "Foz"}, edgesOf(g.Edges[simple+"Foo"]))
	suite.ElementsMatch([]string{simple + "Foz"}, edgesOf(g.Edges[simple+"FooBaz"]))
	suite.ElementsMatch([]string{simple + "Foo", simple + "Foz"}, edgesOf(g.Edges[annotated+"LorenMethod"]))
}

func (suite *NewGraphFromEntriesTestSuite) TestRules() {
	testCases := []struct {
		name    string
		entries []annotation.Entry
		isErr   func(error) bool
		edges   map[string][]string
	}{
		{
			name: "provide and inject on the same function",
			entries: []annotation.Entry{
				fn("NewA", nil, []string{"*A"}, provide("index=0")),
				fn("NewB", []string{"*A"}, []string{"*B"}, provide("index=0"), inject("index=0")),
			},
			edges: map[string][]string{"x/p.NewA": {"x/p.NewB"}},
		},
		{
			name: "unannotated parameter without provider is left to the application",
			entries: []annotation.Entry{
				fn("NewB", []string{"context.Context"}, []string{"*B"}, provide("index=0")),
			},
		},
		{
			name: "unannotated parameter is linked to a default provider",
			entries: []annotation.Entry{
				fn("NewA", nil, []string{"*A"}, provide("index=0")),
				fn("Run", []string{"*A"}, nil, invoke()),
			},
			edges: map[string][]string{"x/p.NewA": {"x/p.Run"}},
		},
		{
			name: "group collects every provider",
			entries: []annotation.Entry{
				fn("One", nil, []string{"string"}, provide("group=names,index=0")),
				fn("Two", nil, []string{"string"}, provide("group=names,index=0")),
				fn("Run", []string{"[]string"}, nil, inject("group=names,index=0"), invoke()),
			},
			edges: map[string][]string{"x/p.One": {"x/p.Run"}, "x/p.Two": {"x/p.Run"}},
		},
		{
			name: "named value is distinct from the default one",
			entries: []annotation.Entry{
				fn("NewA", nil, []string{"*A"}, provide("index=0")),
				fn("NewOtherA", nil, []string{"*A"}, provide("name=other,index=0")),
				fn("Run", []string{"*A"}, nil, inject("name=other,index=0"), invoke()),
			},
			edges: map[string][]string{"x/p.NewOtherA": {"x/p.Run"}},
		},
		{
			name: "duplicate default provider",
			entries: []annotation.Entry{
				fn("NewA", nil, []string{"*A"}, provide("index=0")),
				fn("NewOtherA", nil, []string{"*A"}, provide("index=0")),
			},
			isErr: errors.IsAlreadyExists,
		},
		{
			name: "named inject without provider",
			entries: []annotation.Entry{
				fn("NewA", nil, []string{"*A"}, provide("index=0")),
				fn("Run", []string{"*A"}, nil, inject("name=other,index=0"), invoke()),
			},
			isErr: errors.IsNotFound,
		},
		{
			name: "cycle",
			entries: []annotation.Entry{
				fn("NewA", []string{"*B"}, []string{"*A"}, provide("index=0")),
				fn("NewB", []string{"*A"}, []string{"*B"}, provide("index=0")),
			},
			isErr: errors.IsNotValid,
		},
		{
			name: "self dependency",
			entries: []annotation.Entry{
				fn("Wrap", []string{"*A"}, []string{"*A"}, provide("index=0"), inject("index=0")),
			},
			isErr: errors.IsNotValid,
		},
		{
			name:    "provide index out of range",
			entries: []annotation.Entry{fn("NewA", nil, []string{"*A"}, provide("index=1"))},
			isErr:   errors.IsNotValid,
		},
		{
			name:    "provide index declared twice",
			entries: []annotation.Entry{fn("NewA", nil, []string{"*A"}, provide("index=0"), provide("name=x,index=0"))},
			isErr:   errors.IsNotValid,
		},
		{
			name:    "provide on the error result",
			entries: []annotation.Entry{fn("NewA", nil, []string{"*A", "error"}, provide("index=1"))},
			isErr:   errors.IsNotValid,
		},
		{
			name:    "provide with name and group",
			entries: []annotation.Entry{fn("NewA", nil, []string{"*A"}, provide("name=x,group=y,index=0"))},
			isErr:   errors.IsNotValid,
		},
		{
			name:    "provide and invoke",
			entries: []annotation.Entry{fn("NewA", nil, []string{"*A"}, provide("index=0"), invoke())},
			isErr:   errors.IsNotValid,
		},
		{
			name:    "inject without provide or invoke",
			entries: []annotation.Entry{fn("Use", []string{"*A"}, nil, inject("index=0"))},
			isErr:   errors.IsNotValid,
		},
		{
			name:    "group inject on a non-slice parameter",
			entries: []annotation.Entry{fn("Run", []string{"string"}, nil, inject("group=names,index=0"), invoke())},
			isErr:   errors.IsNotValid,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			g, err := NewGraphFromEntries(context.Background(), tc.entries)
			if tc.isErr != nil {
				suite.Error(err)
				suite.True(tc.isErr(err), "unexpected error type: %v", err)
				return
			}
			suite.Require().NoError(err)

			for from, to := range tc.edges {
				suite.ElementsMatch(to, edgesOf(g.Edges[from]), "edges of %s", from)
			}
			total := 0
			for _, adj := range g.Edges {
				total += len(adj)
			}
			want := 0
			for _, to := range tc.edges {
				want += len(to)
			}
			suite.Equal(want, total, "number of edges")
		})
	}
}

func (suite *NewGraphFromEntriesTestSuite) TestIgnoresEntriesWithoutInjectAnnotations() {
	entries := []annotation.Entry{
		fn("NewA", nil, []string{"*A"}, provide("index=0")),
		fn("Other", nil, nil, annotation.NewAnnotation("MyAnnotation", "")),
		{Path: "x/p", Package: "p", Struct: "S", Annotations: []annotation.Annotation{provide("index=0")}},
	}

	g, err := NewGraphFromEntries(context.Background(), entries)
	suite.Require().NoError(err)
	suite.Equal([]string{"x/p.NewA"}, g.SortedKeys())
}

func (suite *NewGraphFromEntriesTestSuite) TestDeduplicatesEntries() {
	entry := fn("NewA", nil, []string{"*A"}, provide("index=0"))

	g, err := NewGraphFromEntries(context.Background(), []annotation.Entry{entry, entry})
	suite.Require().NoError(err)
	suite.Equal([]string{"x/p.NewA"}, g.SortedKeys())
}

// fn builds the entry of a function in package x/p.
func fn(name string, params, results []string, anns ...annotation.Annotation) annotation.Entry {
	entry := annotation.Entry{
		Path:        "x/p",
		Package:     "p",
		Func:        annotation.EntryFunc{Name: name},
		Annotations: anns,
	}
	for _, tp := range params {
		entry.Func.Parameters = append(entry.Func.Parameters, annotation.EntryFuncType{Type: tp})
	}
	for _, tp := range results {
		entry.Func.Results = append(entry.Func.Results, annotation.EntryFuncType{Type: tp})
	}
	return entry
}

func provide(value string) annotation.Annotation { return annotation.NewAnnotation("Provide", value) }
func inject(value string) annotation.Annotation  { return annotation.NewAnnotation("Inject", value) }
func invoke() annotation.Annotation              { return annotation.NewAnnotation("Invoke", "") }

func edgesOf(adj []*graph.Vertex[Component]) []string {
	keys := make([]string, 0, len(adj))
	for _, v := range adj {
		keys = append(keys, v.Key)
	}
	sort.Strings(keys)
	return keys
}
