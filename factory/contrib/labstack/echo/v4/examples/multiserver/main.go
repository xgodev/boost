package main

import (
	"context"
	"fmt"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"github.com/xgodev/boost/wrapper/config"
	"net/http"
	"time"

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

	ctx1 := context.Background()

	srv1, _ := echo.NewServer(ctx1,
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
