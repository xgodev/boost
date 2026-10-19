package elasticsearch

import (
	"context"
	"testing"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/xgodev/boost/factory"
)

// Each plugin type must stay an alias of factory.Plugin: plugins written
// against either name keep compiling and stay interchangeable.
func TestPluginTypesAreFactoryPlugin(t *testing.T) {
	var _ factory.Plugin[*elasticsearch.Client] = Plugin(func(context.Context, *elasticsearch.Client) error { return nil })
}
