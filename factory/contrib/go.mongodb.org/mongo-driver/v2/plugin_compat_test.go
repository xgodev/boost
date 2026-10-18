package mongo

import (
	"context"
	"testing"

	"github.com/xgodev/boost/factory"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Each plugin type must stay an alias of factory.Plugin: plugins written
// against either name keep compiling and stay interchangeable.
func TestPluginTypesAreFactoryPlugin(t *testing.T) {
	var _ factory.Plugin[*options.ClientOptions] = ClientOptionsPlugin(func(context.Context, *options.ClientOptions) error { return nil })
	var _ factory.Plugin[*mongo.Client] = ClientPlugin(func(context.Context, *mongo.Client) error { return nil })
}
