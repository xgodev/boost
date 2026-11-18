package jwt

import (
	"time"

	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root      = echo.PluginsRoot + ".jwt"
	enabled   = ".enabled"
	jwksURL   = ".jwksUrl"
	issuer    = ".issuer"
	audience  = ".audience"
	skipPaths = ".skipPaths"
	leeway    = ".leeway"
)

func init() {
	ConfigAdd(root)
}

// ConfigAdd registers the jwt plugin keys under path.
func ConfigAdd(path string) {
	config.Add(path+enabled, false, "enable/disable jwt authentication middleware")
	config.Add(path+jwksURL, "", "url of the JWKS holding the public keys that sign the tokens")
	config.Add(path+issuer, "", "issuer (iss) every token must carry")
	config.Add(path+audience, "", "audience (aud) every token must carry")
	config.Add(path+skipPaths, []string{"/health", "/ready", "/live"}, "request paths that need no token")
	config.Add(path+leeway, 30*time.Second, "clock skew tolerated on exp, nbf and iat")
}
