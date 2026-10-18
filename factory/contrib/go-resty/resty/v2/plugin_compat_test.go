package resty

import (
	"context"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/xgodev/boost/factory"
)

// Each plugin type must stay an alias of factory.Plugin: plugins written
// against either name keep compiling and stay interchangeable.
func TestPluginTypesAreFactoryPlugin(t *testing.T) {
	var _ factory.Plugin[*resty.Client] = Plugin(func(context.Context, *resty.Client) error { return nil })
}
