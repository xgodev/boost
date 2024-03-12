package main

import (
	"context"
	"net/http"

	e "github.com/labstack/echo/v4"
	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/factory/labstack/echo.v4"
	"github.com/xgodev/boost/factory/labstack/echo.v4/plugins/contrib/americanas-go/health.v1"
	logplugin "github.com/xgodev/boost/factory/labstack/echo.v4/plugins/contrib/americanas-go/log.v1"
	status "github.com/xgodev/boost/factory/labstack/echo.v4/plugins/contrib/americanas-go/rest-response.v1"
	"github.com/xgodev/boost/factory/labstack/echo.v4/plugins/native/cors"
	"github.com/xgodev/boost/factory/labstack/echo.v4/plugins/native/gzip"
	"github.com/xgodev/boost/factory/labstack/echo.v4/plugins/native/requestid"
	ilog "github.com/xgodev/boost/factory/xgodev/boost.v1/log"
)

const HelloWorldEndpoint = "app.endpoint.helloworld"

func init() {
	config.Add(HelloWorldEndpoint, "/hello-world", "helloworld endpoint")
}

type Config struct {
	App struct {
		Endpoint struct {
			Helloworld string
		}
	}
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

	config.Load()

	ilog.New()

	c := Config{}

	err := config.Unmarshal(&c)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	srv := echo.NewServer(ctx,
		cors.Register,
		requestid.Register,
		gzip.Register,
		logplugin.Register,
		status.Register,
		health.Register)

	srv.GET(c.App.Endpoint.Helloworld, Get)

	srv.Serve(ctx)
}
