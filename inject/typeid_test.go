package inject

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xgodev/boost/annotation"
)

func TestQualifyType(t *testing.T) {
	entry := annotation.Entry{
		Path:    "example.com/app/simple",
		Package: "simple",
		Imports: map[string]string{"annotated": "example.com/lib/annotated", "context": "context"},
	}
	tests := []struct {
		tp   string
		want string
	}{
		{"*Loren", "*example.com/app/simple.Loren"},
		{"Loren", "example.com/app/simple.Loren"},
		{"*annotated.Loren", "*example.com/lib/annotated.Loren"},
		{"*unknown.Loren", "*unknown.Loren"},
		{"string", "string"},
		{"error", "error"},
		{"[]*Loren", "[]*example.com/app/simple.Loren"},
		{"map[string]Loren", "map[string]example.com/app/simple.Loren"},
		{"func(ctx context.Context, l *Loren) error", "func(ctx context.Context, l *example.com/app/simple.Loren) error"},
		{"chan<- Loren", "chan<- example.com/app/simple.Loren"},
		{"not a type(", "not a type("},
	}
	for _, tt := range tests {
		t.Run(tt.tp, func(t *testing.T) {
			assert.Equal(t, tt.want, qualifyType(entry, tt.tp))
		})
	}

	t.Run("two packages with the same name are two types", func(t *testing.T) {
		other := entry
		other.Path = "example.com/app/other/simple"
		assert.NotEqual(t, qualifyType(entry, "*Loren"), qualifyType(other, "*Loren"))
	})
}

func TestQualifyName(t *testing.T) {
	assert.Equal(t, "*simple.Loren", qualifyName("simple", "*Loren"))
	assert.Equal(t, "*annotated.Loren", qualifyName("simple", "*annotated.Loren"))
	assert.Equal(t, "map[string]simple.Loren", qualifyName("simple", "map[string]Loren"))
}

func TestSliceElem(t *testing.T) {
	elem, ok := sliceElem("[]*Loren")
	assert.True(t, ok)
	assert.Equal(t, "*Loren", elem)

	_, ok = sliceElem("[3]Loren")
	assert.False(t, ok)

	_, ok = sliceElem("*Loren")
	assert.False(t, ok)
}

func TestAnnotationTag(t *testing.T) {
	assert.Equal(t, `name:"A"`, (&Annotation{Name: "A"}).Tag())
	assert.Equal(t, `group:"handlers"`, (&Annotation{Group: "handlers"}).Tag())
	assert.Equal(t, "", (&Annotation{}).Tag())
	assert.Equal(t, "named_A", (&Annotation{Name: "A"}).ID())
	assert.Equal(t, "grouped_handlers", (&Annotation{Group: "handlers"}).ID())
	assert.Equal(t, "default", (&Annotation{}).ID())
}

func TestWithImports(t *testing.T) {
	shared := map[string]string{"context": "context"}
	consumer := annotation.Entry{
		Path: "example.com/app/app", Package: "app", Imports: shared,
		Func: annotation.EntryFunc{Parameters: []annotation.EntryFuncType{
			{Type: "*quote.Client"}, {Type: "context.Context"}, {Type: "*boost.Options"}, {Type: "*nowhere.T"}, {Type: "("},
		}, Results: []annotation.EntryFuncType{{Type: "map[string]*quote.Client"}}},
	}
	entries := []annotation.Entry{
		consumer,
		{Path: "example.com/app/quote", Package: "quote"},
		{Path: "example.com/app/a/boost", Package: "boost"},
		{Path: "example.com/app/b/boost", Package: "boost"},
	}

	got := withImports(entries)

	assert.Equal(t, map[string]string{"context": "context", "quote": "example.com/app/quote"}, got[0].Imports,
		"a qualifier left out resolves to the one package with that name; two candidates or none stay unresolved")
	assert.Equal(t, map[string]string{"context": "context"}, shared, "the map the file's entries share is not changed")
	assert.Nil(t, got[1].Imports, "an entry with nothing to resolve is left as is")
}
