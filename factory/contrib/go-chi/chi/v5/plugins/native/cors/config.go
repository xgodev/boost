package cors

import (
	"net/http"

	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root               = chi.PluginsRoot + ".cors"
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

// ConfigAdd registers the cors middleware plugin config keys under path.
func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable cors middleware")
	config.Add(path+allowedOrigins, []string{"*"}, "cors allowed origins")
	config.Add(path+allowedHeaders, []string{"Origin", "Content-Type", "Accept", "Authorization"}, "cors allowed headers")
	config.Add(path+allowedMethods, []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete}, "cors allowed methods")
	config.Add(path+allowedCredentials, false, "cors allowed credentials")
	config.Add(path+exposedHeaders, []string{}, "cors exposed headers")
	config.Add(path+maxAge, 5200, "cors max age (seconds)")
}
