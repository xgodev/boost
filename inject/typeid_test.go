package inject

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQualifyType(t *testing.T) {
	tests := []struct {
		tp   string
		want string
	}{
		{"*Loren", "*simple.Loren"},
		{"Loren", "simple.Loren"},
		{"*annotated.Loren", "*annotated.Loren"},
		{"string", "string"},
		{"error", "error"},
		{"[]*Loren", "[]*simple.Loren"},
		{"map[string]Loren", "map[string]simple.Loren"},
		{"func(ctx context.Context, l *Loren) error", "func(ctx context.Context, l *simple.Loren) error"},
		{"chan<- Loren", "chan<- simple.Loren"},
	}
	for _, tt := range tests {
		t.Run(tt.tp, func(t *testing.T) {
			assert.Equal(t, tt.want, qualifyType("simple", tt.tp))
		})
	}
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
