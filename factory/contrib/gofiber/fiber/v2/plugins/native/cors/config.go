package cors

import (
	"net/http"

	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root               = fiber.PluginsRoot + ".cors"
	enabled            = ".enabled"
	allowedOrigins     = ".allowed.origins"
	allowedHeaders     = ".allowed.headers"
	allowedMethods     = ".allowed.methods"
	allowedCredentials = ".allowed.credentials"
	exposedHeaders     = ".exposed.headers"
	maxAge             = ".maxAge"
)

func init() {
	ConfigAdd(root)
}

// ConfigAdd registers the cors middleware plugin config keys under path.
func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable cors middleware")
	config.Add(path+allowedOrigins, []string{"*"}, "cors allowed origins")
	config.Add(path+allowedHeaders, []string{"Origin", "Content-Type", "Accept", "Authorization"}, "cors allowed headers")
	config.Add(path+allowedMethods, []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete}, "cors allowed methods")
	config.Add(path+allowedCredentials, false, "cors allowed credentials, fiber refuses it with the origin *")
	config.Add(path+exposedHeaders, []string{}, "cors exposed headers")
	config.Add(path+maxAge, 5200, "cors max age (seconds)")
}
