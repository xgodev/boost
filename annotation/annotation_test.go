package annotation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAnnotation_Map(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  map[string]interface{}
	}{
		{"pairs", "method=GET, path=/a", map[string]interface{}{"method": "GET", "path": "/a"}},
		{"equals inside the value", "filter=a=b", map[string]interface{}{"filter": "a=b"}},
		{"quoted value keeps commas", `tags="a,b", n=1`, map[string]interface{}{"tags": "a,b", "n": 1}},
		{"parens keep commas", "path=/a(b,c)", map[string]interface{}{"path": "/a(b,c)"}},
		{"typed values", "flag=T, on=true, off=FALSE, ratio=1.5, n=2", map[string]interface{}{
			"flag": "T", "on": true, "off": false, "ratio": 1.5, "n": 2,
		}},
		{"empty", "", map[string]interface{}{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAnnotation("A", tt.value)
			assert.Equal(t, tt.want, got.Map)
		})
	}
}

func TestNewAnnotation_RawValueKeepsTheSourceText(t *testing.T) {
	a := NewAnnotation("Rest", "path=/a, method=GET")

	assert.Equal(t, "path=/a, method=GET", a.RawValue())
}

func TestAnnotation_Decode(t *testing.T) {
	var got struct {
		Name  string `attr:"name"`
		Index int    `attr:"index"`
	}
	a := NewAnnotation("Inject", "name=svc, index=2")

	require.NoError(t, a.Decode(&got))
	assert.Equal(t, "svc", got.Name)
	assert.Equal(t, 2, got.Index)
}

func TestCollector_extractAnnotation(t *testing.T) {
	tests := []struct {
		name      string
		filters   []string
		comment   string
		ok        bool
		wantName  string
		wantValue string
	}{
		{"with values", nil, "// @Rest(method=GET)", true, "Rest", "method=GET"},
		{"without values", nil, "// @Boost", true, "Boost", ""},
		{"trailing space", nil, "// @Boost ", true, "Boost", ""},
		{"trailing space after values", nil, "// @Rest(method=GET) ", true, "Rest", "method=GET"},
		{"nested parens", nil, "// @Rest(path=/a(b))", true, "Rest", "path=/a(b)"},
		{"not an annotation", nil, "// plain comment", false, "", ""},
		{"filter matches", []string{"Inject"}, "// @Inject(name=a)", true, "Inject", "name=a"},
		{"filter matches a name prefix", []string{"Rest"}, "// @RestRouter(path=/)", true, "RestRouter", "path=/"},
		{"filter rejects other names", []string{"Rest"}, "// @Other(a=1)", false, "", ""},
		{"empty parens", nil, "// @Boost()", true, "Boost", ""},
		{"value without key", nil, "// @Rest(GET)", false, "", ""},
		{"unbalanced quote", nil, `// @Rest(a="b)`, false, "", ""},
		{"words after the name", nil, "// @Param query foo bool", false, "", ""},
		{"filter ignores mentions in values", []string{"Inject"}, `// @Rest(desc="uses @Inject")`, false, "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collector{filters: tt.filters}

			got, ok := c.extractAnnotation(tt.comment)

			require.Equal(t, tt.ok, ok)
			if ok {
				assert.Equal(t, tt.wantName, got.Name)
				assert.Equal(t, tt.wantValue, got.RawValue())
			}
		})
	}
}
