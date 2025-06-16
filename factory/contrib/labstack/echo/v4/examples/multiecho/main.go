package main

import (
	"context"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"github.com/xgodev/boost/wrapper/config"
	"net/http"
	"os"

	e "github.com/labstack/echo/v4"
	"github.com/xgodev/boost/extra/multiserver"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/local/extra/health"
	mserver "github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/local/extra/multi-server"
	status "github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/local/model/restresponse"
	logplugin "github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/local/wrapper/log"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/native/cors"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/native/gzip"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/native/requestid"
)

func init() {
	echo.ConfigAdd("boost.factory.echo2")
}

type Response struct {
	Message string
}

func Get(c e.Context) (err error) {

	resp := Response{
		Message: "Hello World!!",
	}

	err = config.Unmarshal(&resp)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func main() {

	os.Setenv("BOOST_FACTORY_ECHO2_PORT", "8086")

	boost.Start()

	ctx1 := context.Background()

	srv1, _ := echo.NewServer(ctx1,
		cors.Register,
		requestid.Register,
		gzip.Register,
		logplugin.Register,
		status.Register,
		health.Register)

	srv1.GET("/test", Get)

	ctx2 := context.Background()

	options2, err := echo.NewOptionsWithPath("boost.factory.echo2")
	if err != nil {
		panic(err)
	}

	srv2 := echo.NewServerWithOptions(ctx2, options2,
		cors.Register,
		requestid.Register,
		gzip.Register,
		logplugin.Register,
		status.Register,
		mserver.Register,
		health.Register)

	srv2.GET("/test", Get)

	multiserver.Serve(context.Background(), srv1, srv2)
}
