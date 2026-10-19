package echo

import (
	"context"
	"testing"

	"github.com/xgodev/boost/factory"
)

// Each plugin type must stay an alias of factory.Plugin: plugins written
// against either name keep compiling and stay interchangeable.
func TestPluginTypesAreFactoryPlugin(t *testing.T) {
	var _ factory.Plugin[*Server] = Plugin(func(context.Context, *Server) error { return nil })
}
