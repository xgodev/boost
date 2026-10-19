// Package rest holds what the rest server of app api needs besides its
// handlers.
package rest

import (
	boostecho "github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/extra/error_handler"
)

// Plugins turn boost errors into HTTP statuses: a NotFound answers 404.
// @Provide(index=0)
func Plugins() []boostecho.Plugin {
	return []boostecho.Plugin{error_handler.Register}
}
