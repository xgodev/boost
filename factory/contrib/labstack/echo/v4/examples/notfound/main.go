package main

import (
	"context"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"net/http"
	"os"

	e "github.com/labstack/echo/v4"
	prometheus "github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/contrib/prometheus/client_golang/v1"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/extra/error_handler"
	logplugin "github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/local/wrapper/log"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/native/cors"
	"github.com/xgodev/boost/factory/local/wrapper/log"
	"github.com/xgodev/boost/model/errors"
)

func errorHandler(c e.Context) (err error) {
	return errors.NotFoundf("example")
}

func helloHandler(c e.Context) (err error) {
	c.String(http.StatusOK, "hello world")
	return nil
}

func main() {

	os.Setenv("BOOST_FACTORY_LOGRUS_CONSOLE_LEVEL", "TRACE")

	boost.Start()
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
