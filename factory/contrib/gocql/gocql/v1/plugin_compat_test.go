package gocql

import (
	"context"
	"testing"

	"github.com/gocql/gocql"
	"github.com/xgodev/boost/factory"
)

// Each plugin type must stay an alias of factory.Plugin: plugins written
// against either name keep compiling and stay interchangeable.
func TestPluginTypesAreFactoryPlugin(t *testing.T) {
	var _ factory.Plugin[*gocql.Session] = Plugin(func(context.Context, *gocql.Session) error { return nil })
}
