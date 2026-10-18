package inject

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost/annotation"
	"github.com/xgodev/boost/model/errors"
)

func deflt() annotation.Annotation { return annotation.NewAnnotation("Default", "") }

// impl builds New in package stock/<pkg>, an implementation of stock.Service.
func impl(pkg string, anns ...annotation.Annotation) annotation.Entry {
	entry := at("stock/"+pkg, "New", nil, []string{"stock.Service"}, anns...)
	entry.Imports = map[string]string{"stock": testModule + "/stock"}
	return entry
}

// consumer builds NewService in package order, which takes a stock.Service.
func consumer(anns ...annotation.Annotation) annotation.Entry {
	entry := at("order", "NewService", []string{"stock.Service"}, []string{"*Service"}, append(anns, provide("index=0"))...)
	entry.Imports = map[string]string{"stock": testModule + "/stock"}
	return entry
}

func stockEntries() []annotation.Entry {
	return []annotation.Entry{
		impl("memory", provide("index=0"), deflt()),
		impl("unlimited", provide("index=0")),
		consumer(),
	}
}

func TestImplementations(t *testing.T) {
	sets, err := Implementations(newTestGraph(t, stockEntries()))
	require.NoError(t, err)
	require.Len(t, sets, 1)

	set := sets[0]
	assert.Equal(t, "Service", set.Type)
	assert.Equal(t, testModule+"/stock", set.Path)
	assert.Equal(t, "", set.Name)
	assert.Equal(t, "memory", set.Default)
	assert.Equal(t, []string{"memory", "unlimited"}, set.Names())
	assert.Equal(t, testModule+"/stock/memory.New", set.Members["memory"].Key())
	assert.Equal(t, "Service", set.Ident())
	assert.Equal(t, "stock.implementation", set.Key(testModule))
	assert.Equal(t, "implementation/stock/service_implementation.go", set.File("example.com/app/gen"))

	g := newTestGraph(t, stockEntries())
	member, ok := ImplementationOf(g, g.Vertices[testModule+"/stock/unlimited.New"].Value)
	require.True(t, ok)
	assert.Equal(t, set.Key(testModule), member.Key(testModule))
	_, ok = ImplementationOf(g, g.Vertices[testModule+"/order.NewService"].Value)
	assert.False(t, ok, "a consumer is not a member")
}

func TestImplementations_Named(t *testing.T) {
	sets, err := Implementations(newTestGraph(t, []annotation.Entry{
		impl("memory", provide("name=formal,index=0"), deflt()),
		impl("unlimited", provide("name=formal,index=0")),
	}))
	require.NoError(t, err)
	require.Len(t, sets, 1)
	assert.Equal(t, "formal", sets[0].Name)
	assert.Equal(t, "ServiceFormal", sets[0].Ident())
	assert.Equal(t, "stock.formal.implementation", sets[0].Key(testModule))
	assert.Equal(t, "implementation/stock/serviceformal_implementation.go", sets[0].File("example.com/app/gen"))
}

func TestImplementation_Key(t *testing.T) {
	testCases := []struct {
		name string
		set  Implementation
		want string
	}{
		{name: "package in the module", set: Implementation{Path: testModule + "/stock"}, want: "stock.implementation"},
		{name: "nested package", set: Implementation{Path: testModule + "/internal/domain/quote"}, want: "internal.domain.quote.implementation"},
		{name: "named value", set: Implementation{Path: testModule + "/stock", Name: "formal"}, want: "stock.formal.implementation"},
		{name: "module root", set: Implementation{Path: testModule}, want: "app.implementation"},
		{name: "outside the module", set: Implementation{Path: "github.com/cloudevents/sdk-go/v2"}, want: "github.com.cloudevents.sdkGo.v2.implementation"},
		{name: "underscore splits", set: Implementation{Path: testModule + "/my_pkg/x"}, want: "my.pkg.x.implementation"},
		{name: "upper case is lowered", set: Implementation{Path: testModule + "/Stock", Name: "Formal"}, want: "stock.formal.implementation"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.set.Key(testModule))
		})
	}
}

func TestImplementations_Rules(t *testing.T) {
	other := func(pkg string, anns ...annotation.Annotation) annotation.Entry {
		entry := impl(pkg, anns...)
		entry.Path = "example.com/other/" + pkg
		return entry
	}

	testCases := []struct {
		name    string
		entries []annotation.Entry
		isErr   func(error) bool
		message string
	}{
		{
			name:    "two providers without @Default",
			entries: []annotation.Entry{impl("memory", provide("index=0")), impl("unlimited", provide("index=0"))},
			isErr:   errors.IsAlreadyExists,
			message: "@Default",
		},
		{
			name:    "two @Default",
			entries: []annotation.Entry{impl("memory", provide("index=0"), deflt()), impl("unlimited", provide("index=0"), deflt())},
			isErr:   errors.IsNotValid,
			message: "more than one @Default",
		},
		{
			name:    "@Default without another implementation",
			entries: []annotation.Entry{impl("memory", provide("index=0"), deflt())},
			isErr:   errors.IsNotValid,
			message: "no other implementation",
		},
		{
			name:    "@Default without @Provide",
			entries: []annotation.Entry{impl("memory", deflt())},
			isErr:   errors.IsNotValid,
		},
		{
			name:    "@Default with @Invoke",
			entries: []annotation.Entry{at("app", "Run", nil, nil, invoke(), deflt())},
			isErr:   errors.IsNotValid,
		},
		{
			name: "@Default on a group",
			entries: []annotation.Entry{
				impl("memory", provide("group=stocks,index=0"), deflt()),
				impl("unlimited", provide("group=stocks,index=0")),
			},
			isErr: errors.IsNotValid,
		},
		{
			name:    "@Default with arguments",
			entries: []annotation.Entry{impl("memory", provide("index=0"), annotation.NewAnnotation("Default", "x=1")), impl("unlimited", provide("index=0"))},
			isErr:   errors.IsNotValid,
		},
		{
			name: "a member with two @Provide",
			entries: []annotation.Entry{
				impl("memory", provide("index=0"), deflt()),
				func() annotation.Entry {
					e := impl("unlimited", provide("index=0"), provide("index=1"))
					e.Func.Results = append(e.Func.Results, annotation.EntryFuncType{Type: "*Clock"})
					return e
				}(),
			},
			isErr:   errors.IsNotValid,
			message: "exactly one @Provide",
		},
		{
			name:    "two members in packages with the same name",
			entries: []annotation.Entry{impl("memory", provide("index=0"), deflt()), other("memory", provide("index=0"))},
			isErr:   errors.IsNotValid,
			message: "memory",
		},
		{
			name: "members resolve the type to different packages",
			entries: []annotation.Entry{
				impl("memory", provide("index=0"), deflt()),
				func() annotation.Entry {
					e := impl("unlimited", provide("index=0"))
					e.Imports = map[string]string{"stock": "example.com/vendor/stock"}
					return e
				}(),
			},
			isErr: errors.IsNotValid,
		},
		{
			name: "unnamed type",
			entries: []annotation.Entry{
				at("one", "Names", nil, []string{"[]string"}, provide("index=0"), deflt()),
				at("two", "Names", nil, []string{"[]string"}, provide("index=0")),
			},
			isErr:   errors.IsNotValid,
			message: "named type",
		},
		{
			name: "consumer in a member's package",
			entries: []annotation.Entry{
				impl("memory", provide("index=0"), deflt()),
				impl("unlimited", provide("index=0")),
				func() annotation.Entry {
					e := at("stock/memory", "NewReport", []string{"stock.Service"}, []string{"*Report"}, provide("index=0"))
					e.Imports = map[string]string{"stock": testModule + "/stock"}
					return e
				}(),
			},
			isErr:   errors.IsNotValid,
			message: "same package",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewGraphFromEntries(context.Background(), tc.entries)
			require.Error(t, err)
			assert.True(t, tc.isErr(err), "unexpected error type: %v", err)
			assert.Contains(t, err.Error(), tc.message)
		})
	}
}

func TestImplementations_LinksEveryMember(t *testing.T) {
	g := newTestGraph(t, stockEntries())
	for _, member := range []string{"memory", "unlimited"} {
		assert.Equal(t, []string{testModule + "/order.NewService"}, edgesOf(g.Edges[testModule+"/stock/"+member+".New"]),
			"cycles and reachability see every implementation")
	}
}
