package main

import (
	"context"

	"github.com/wesovilabs/beyond/api"
	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/local/extra/health"
	status "github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/local/model/restresponse"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/local/wrapper/log"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/native/cors"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/native/gzip"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/native/requestid"
	ilog "github.com/xgodev/boost/factory/local/wrapper/log"
)

const Endpoint = "app.endpoint.google"

func init() {
	config.Add(Endpoint, "/google", "google endpoint")
}

func Beyond() *api.Beyond {
	return api.New().
		WithBefore(NewTracingAdvice, "handler.Get(...)").
		WithBefore(NewTracingAdviceWithPrefix("[beyond]"), "handler.*(...)...")
}

func main() {

	var err error

	config.Load()

	c := Config{}

	err = config.Unmarshal(&c)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	ilog.New()

	srv := echo.NewServer(ctx,
		cors.Register,
		requestid.Register,
		gzip.Register,
		log.Register,
		status.Register,
		health.Register)

	srv.GET(c.App.Endpoint.Google, Get)

	srv.Serve(ctx)
}
