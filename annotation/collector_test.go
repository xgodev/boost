package annotation

import (
	"testing"
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
