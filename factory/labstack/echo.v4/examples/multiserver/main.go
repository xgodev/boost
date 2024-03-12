package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	e "github.com/labstack/echo/v4"
	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/factory/labstack/echo.v4"
	"github.com/xgodev/boost/factory/labstack/echo.v4/plugins/contrib/americanas-go/health.v1"
	logplugin "github.com/xgodev/boost/factory/labstack/echo.v4/plugins/contrib/americanas-go/log.v1"
	mserver "github.com/xgodev/boost/factory/labstack/echo.v4/plugins/contrib/americanas-go/multi-server.v1"
	status "github.com/xgodev/boost/factory/labstack/echo.v4/plugins/contrib/americanas-go/rest-response.v1"
	"github.com/xgodev/boost/factory/labstack/echo.v4/plugins/native/cors"
	"github.com/xgodev/boost/factory/labstack/echo.v4/plugins/native/gzip"
	"github.com/xgodev/boost/factory/labstack/echo.v4/plugins/native/requestid"
	ilog "github.com/xgodev/boost/factory/xgodev/boost.v1/log"
	"github.com/xgodev/boost/multiserver"
)

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

	ctx1 := context.Background()

	srv1 := echo.NewServer(ctx1,
		cors.Register,
		requestid.Register,
		gzip.Register,
		logplugin.Register,
		status.Register,
		mserver.Register,
		health.Register)

	srv1.GET("/test", Get)

	multiserver.Serve(context.Background(), srv1, &LocalServer{})
}

type LocalServer struct {
}

func (s *LocalServer) Serve(ctx context.Context) {
	time.Sleep(10 * time.Second)
	fmt.Printf("finished")
}

func (s *LocalServer) Shutdown(ctx context.Context) {
}
