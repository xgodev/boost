package annotation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollector_isValidAnnotation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Test1", "// @RestResponse(code=201, type=github.com/xgodev/boost/inject/examples/simple.Response, description=tiam sed efficitur purus at lacinia magna)", true},
		{"Test2", "// @RestRequest(code=201)", true},
		{"Test3", "// @RestAction (type=action, data=123)", true},
		{"Test4", "// @RestAction (type=action, data=123)", true},
		{"Test5", "// @RestAction ( type=action, data=123)", true},
		{"Test6", "// @RestAction ( type=action, data=123 )", true},
		{"Test7", "// @RestAction (type=action,data=123, xpto=456)", true},
		{"Test8", "// @RestRequest(code=201)", true},
		{"Test9", "//@RestAction (type=action, data=123)", false},
		{"Test10", "//@RestAction (type=action, data=123)", false},
		{"Test11", "//@RestAction ( type=action, data=123)", false},
		{"Test12", "//@RestAction ( type=action, data=123 )", false},
		{"Test13", "//@RestAction (type=action,data=123, xpto=456)", false},
		{"Test14", "// FooFunc Lorem ipsum dolor sit amet, consectetur adipiscing elit", false},
		{"Test15", "//FooFunc Lorem ipsum dolor sit amet, consectetur adipiscing elit", false},
		{"Test16", "// @MyAnnotation(code=201)", true},
		{"Test16", "// @MyAnnotation(code=201", false},
		{"Test16", "// @A Param query foo bool true tiam sed efficitur purus", false},
		{"Test16", "// @Invoke", true},
		{"bare names", "// @Middleware(log, hystrix)", true},
		{"bare names and pairs", "// @Middleware(log, name=x)", true},
		{"words are not a name", "// @Middleware(log hystrix)", false},
		{"a name starts with a letter", "// @Middleware(0)", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			c := &Collector{
				filters: []string{""},
			}

			if got := c.isValidAnnotation(tt.input); got != tt.expected {
				t.Errorf("isValidAnnotation() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCollector_PackageNameFallsBackToThePath(t *testing.T) {
	c := &Collector{names: map[string]string{"example.com/sdk-go/v2": "cloudevents"}}

	assert.Equal(t, "cloudevents", c.packageName("example.com/sdk-go/v2"), "the loaded name wins")
	assert.Equal(t, "echo", c.packageName("github.com/labstack/echo/v4"), "major version suffixes are skipped")
	assert.Equal(t, "v4", c.packageName("v4"))
	assert.Equal(t, "strings", c.packageName("strings"))
}
