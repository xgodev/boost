package cors

import (
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v5"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root               = echo.PluginsRoot + ".cors"
	enabled            = ".enabled"
	allowedRoot        = ".allowed"
	allowedOrigins     = allowedRoot + ".origins"
	allowedHeaders     = allowedRoot + ".headers"
	allowedMethods     = allowedRoot + ".methods"
	allowedCredentials = allowedRoot + ".credentials"
	exposedHeaders     = ".exposed.headers"
	maxAge             = ".maxAge"
)

func init() {
	ConfigAdd(root)
}

// ConfigAdd registers the cors plugin config keys under path.
func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable cors middleware")
	config.Add(path+allowedOrigins, []string{"*"}, "cors allowed origins")
	config.Add(path+allowedHeaders, []string{"*"}, "cors allowed headers")
	config.Add(path+allowedMethods, []string{"*"}, "cors allowed methods")
	config.Add(path+allowedCredentials, false, "cors allowed credentials; echo v5 refuses it with the \"*\" origin")
	config.Add(path+exposedHeaders, []string{}, "cors exposed headers")
	config.Add(path+maxAge, 5200, "cors max age (seconds)")
}
