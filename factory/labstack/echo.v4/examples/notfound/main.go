package main

import (
	"context"
	"net/http"
	"os"

	e "github.com/labstack/echo/v4"
	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/errors"
	"github.com/xgodev/boost/factory/labstack/echo.v4"
	logplugin "github.com/xgodev/boost/factory/labstack/echo.v4/plugins/contrib/americanas-go/log.v1"
	prometheus "github.com/xgodev/boost/factory/labstack/echo.v4/plugins/contrib/prometheus/client_golang.v1"
	"github.com/xgodev/boost/factory/labstack/echo.v4/plugins/extra/error_handler"
	"github.com/xgodev/boost/factory/labstack/echo.v4/plugins/native/cors"
	"github.com/xgodev/boost/factory/xgodev/boost.v1/log"
)

func errorHandler(c e.Context) (err error) {
	return errors.NotFoundf("example")
}

func helloHandler(c e.Context) (err error) {
	c.String(http.StatusOK, "hello world")
	return nil
}

func main() {

	os.Setenv("IGNITE_LOGRUS_CONSOLE_LEVEL", "TRACE")

	config.Load()
	log.New()

	ctx := context.Background()

	srv := echo.NewServer(ctx,
		cors.Register,
		logplugin.Register,
		prometheus.Register,
		error_handler.Register)

	srv.GET("/not-found", errorHandler)
	srv.GET("/hello", helloHandler)

	srv.Serve(ctx)
}
