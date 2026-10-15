package annotation

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fixtureModule writes a self-contained Go module under a hidden ancestor
// directory, so the tests prove that the location of the project on disk
// does not change what gets collected.
func fixtureModule(t *testing.T) string {
	t.Helper()

	root := filepath.Join(t.TempDir(), ".hidden-ancestor", "fixture")
	files := map[string]string{
		"go.mod": "module example.com/fixture\n\ngo 1.21\n",
		"app/service.go": `package app

import "example.com/fixture/dep"

// Service handles requests.
// It is the main entry point.
// @Inject(name=svc)
type Service struct{}

// Get fetches a resource.
// @Rest(method=GET, path=/a)
func (s *Service) Get() string { return dep.Hello() }

// @Rest(method=POST)
func (s Service) Post() {}

// @Other
func NotFiltered() {}

// @Boost
func NewService() *Service { return &Service{} }

// @Inject
type Repo interface{ Find() }

type (
	// @Inject(kind=grouped)
	Grouped struct{}
)

type Box[T any] struct{ v T }

// @Rest(method=GET)
func (b *Box[T]) Read() T { return b.v }
`,
		"app/other.go": `package app

// @Rest(method=PUT)
func (s *Service) Put() {}
`,
		"dep/dep.go": `package dep

// @Boost(kind=dep)
func Hello() string { return "hi" }
`,
		"app/.skip/skip.go": `package skip

// @Boost(kind=hidden)
func Hidden() {}
`,
	}
	for name, content := range files {
		path := filepath.Join(root, name)
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
		require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	}
	return root
}

func collectEntries(t *testing.T, opts ...Option) []Entry {
	t.Helper()
	c, err := Collect(append([]Option{WithFilters("Inject", "Rest", "Boost")}, opts...)...)
	require.NoError(t, err)
	return c.Entries()
}

// find returns the entries declared as owner.name ("" owner for package
// functions and "" name for type declarations).
func find(entries []Entry, owner, name string) []Entry {
	var found []Entry
	for _, e := range entries {
		if e.Struct == owner && e.Func.Name == name {
			found = append(found, e)
		}
	}
	return found
}

func TestCollect_FromProjectUnderHiddenAncestor(t *testing.T) {
	root := fixtureModule(t)

	entries := collectEntries(t, WithPath(root))

	svc := find(entries, "Service", "")
	require.Len(t, svc, 1)
	assert.True(t, svc[0].IsStruct())
	assert.Equal(t, "example.com/fixture/app", svc[0].Path)
	assert.Equal(t, "app", svc[0].Package)
	assert.Equal(t, "example.com/fixture", svc[0].Module)
	assert.Equal(t, "service.go", svc[0].File)
	assert.Equal(t, "handles requests.", svc[0].Header.Title)
	assert.Equal(t, "It is the main entry point.", svc[0].Header.Description)
	require.Len(t, svc[0].Annotations, 1)
	assert.Equal(t, "Inject", svc[0].Annotations[0].Name)

	hello := find(entries, "", "Hello")
	require.Len(t, hello, 1, "a package reached by the walk and by an import is collected once")
	assert.Equal(t, "example.com/fixture/dep", hello[0].Path)

	assert.Empty(t, find(entries, "", "Hidden"), "hidden directories below the base path are skipped")
}

func TestCollect_ClassifiesMethods(t *testing.T) {
	entries := collectEntries(t, WithPath(fixtureModule(t)))

	for _, method := range []string{"Get", "Post", "Put"} {
		got := find(entries, "Service", method)
		require.Len(t, got, 1, method)
		assert.True(t, got[0].IsMethod(), method)
	}
	assert.Equal(t, "other.go", find(entries, "Service", "Put")[0].File)

	read := find(entries, "Box", "Read")
	require.Len(t, read, 1, "generic receivers are methods")
	assert.True(t, read[0].IsMethod())

	for _, name := range []string{"Post", "Put", "Read"} {
		assert.Empty(t, find(entries, "", name), "%s must not be reported as a package function", name)
	}

	newSvc := find(entries, "", "NewService")
	require.Len(t, newSvc, 1)
	assert.True(t, newSvc[0].IsFunc())
}

func TestCollect_TypeDeclarations(t *testing.T) {
	entries := collectEntries(t, WithPath(fixtureModule(t)))

	grouped := find(entries, "Grouped", "")
	require.Len(t, grouped, 1, "a grouped type declaration keeps its own doc comment")
	assert.Equal(t, "grouped", grouped[0].Annotations[0].Map["kind"])

	assert.Empty(t, find(entries, "Repo", ""), "interfaces are not structs")
}

func TestCollect_FiltersMatchTheAnnotationName(t *testing.T) {
	entries := collectEntries(t, WithPath(fixtureModule(t)))

	assert.Empty(t, find(entries, "", "NotFiltered"), "@Other does not match the filters")
}

func TestCollect_WithPackagesMatchesImportPaths(t *testing.T) {
	root := fixtureModule(t)

	entries := collectEntries(t, WithPath(root), WithPackages("example.com/fixture/dep"))

	require.Len(t, entries, 1)
	assert.Equal(t, "Hello", entries[0].Func.Name)
}

func TestCollect_FollowsImportsOutsideThePath(t *testing.T) {
	root := fixtureModule(t)

	entries := collectEntries(t, WithPath(filepath.Join(root, "app")), WithPackages("example.com/fixture"))

	assert.Len(t, find(entries, "", "Hello"), 1, "imported packages allowed by WithPackages are collected")
}

func TestCollect_RequiresPath(t *testing.T) {
	_, err := Collect()
	assert.Error(t, err)
}

func TestCollector_Queries(t *testing.T) {
	ctx := EntryFuncType{Name: "ctx", Type: "context.Context"}
	errType := EntryFuncType{Type: "error"}
	c := &Collector{entries: []Entry{
		{Struct: "Service", Annotations: []Annotation{{Name: "Inject"}}},
		{Func: EntryFunc{Name: "Handle", Parameters: []EntryFuncType{ctx}, Results: []EntryFuncType{errType}},
			Annotations: []Annotation{{Name: "RestRouter"}}},
		{Func: EntryFunc{Name: "Provide", Results: []EntryFuncType{{Type: "*Service"}}},
			Annotations: []Annotation{{Name: "Provide"}}},
	}}

	names := func(entries []Entry) (out []string) {
		for _, e := range entries {
			out = append(out, e.Struct+e.Func.Name)
		}
		return out
	}

	assert.Equal(t, []string{"Service"}, names(c.EntriesWith("Inject")))
	assert.Equal(t, []string{"Handle"}, names(c.EntriesWithPrefix("Rest")))
	assert.Equal(t, []string{"Handle"}, names(c.EntriesWithArgType("RestRouter", "context.Context")))
	assert.Empty(t, c.EntriesWithArgType("Provide", "context.Context"))
	assert.Equal(t, []string{"Provide"}, names(c.EntriesWithResultType("Provide", "*Service")))
	assert.Empty(t, c.EntriesWithResultType("RestRouter", "*Service"))
}
