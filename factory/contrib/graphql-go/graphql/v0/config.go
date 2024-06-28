package graphql

import (
	"github.com/graphql-go/handler"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root             = "boost.factory.graphql"
	handlerConfig    = root + ".handler"
	pretty           = handlerConfig + ".pretty"
	enableGraphiQL   = handlerConfig + ".graphiQL"
	enablePlayground = handlerConfig + ".playground"
)

func init() {
	config.Add(pretty, false, "enable/disable pretty print")
	config.Add(enableGraphiQL, false, "enable/disable GraphiQL")
	config.Add(enablePlayground, true, "enable/disable Playground")
}

// DefaultHandlerConfig unmarshals the default graphql handler config.
func DefaultHandlerConfig() (*handler.Config, error) {
	return boost.NewOptionsWithPath[handler.Config](handlerConfig)
}
