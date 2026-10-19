// Package rest holds what the rest server of app api needs besides its
// handlers.
package rest

import (
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/extra/error_handler"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/local/wrapper/log"
)

// Plugins log the requests and answer errors as JSON.
// @Provide(index=0)
func Plugins() []boostecho.Plugin {
	return []boostecho.Plugin{
		log.Register,
		error_handler.Register,
	}
}
