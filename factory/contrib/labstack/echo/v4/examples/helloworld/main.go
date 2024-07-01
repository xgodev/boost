package main

import (
	"context"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/wrapper/config"
	"net/http"

	e "github.com/labstack/echo/v4"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/local/extra/health"
	status "github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/local/model/restresponse"
	logplugin "github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/local/wrapper/log"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/native/cors"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/native/gzip"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/native/requestid"
	ilog "github.com/xgodev/boost/factory/local/wrapper/log"
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

	boost.Start()

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
