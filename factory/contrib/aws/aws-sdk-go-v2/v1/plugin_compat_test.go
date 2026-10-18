package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/xgodev/boost/factory"
)

// Each plugin type must stay an alias of factory.Plugin: plugins written
// against either name keep compiling and stay interchangeable.
func TestPluginTypesAreFactoryPlugin(t *testing.T) {
	var _ factory.Plugin[*aws.Config] = Plugin(func(context.Context, *aws.Config) error { return nil })
}
