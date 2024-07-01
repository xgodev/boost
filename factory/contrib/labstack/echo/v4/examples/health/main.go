package main

import (
	"context"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/extra/health"
	"net/http"
	"os"

	e "github.com/labstack/echo/v4"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	h "github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/local/extra/health"
)

func helloHandler(c e.Context) (err error) {
	c.String(http.StatusOK, "hello world")
	return nil
}

type MyChecker struct {
}

func (c *MyChecker) Check(ctx context.Context) error {
	return nil
}

func main() {

	os.Setenv("BOOST_FACTORY_LOGRUS_CONSOLE_LEVEL", "TRACE")

	boost.Start()

	ctx := context.Background()

	hc := health.NewHealthChecker("teste", "teste", &MyChecker{}, true, true)
	health.Add(hc)

	srv := echo.NewServer(ctx, h.Register)

	srv.GET("/hello", helloHandler)

	srv.Serve(ctx)
}
