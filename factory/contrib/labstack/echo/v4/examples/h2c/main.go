package main

import (
	"context"
	"fmt"
	"github.com/xgodev/boost"
	"net/http"
	"os"

	e "github.com/labstack/echo/v4"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4"
	"github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/extra/error_handler"
	logplugin "github.com/xgodev/boost/factory/contrib/labstack/echo/v4/plugins/local/wrapper/log"
)

func init() {
	os.Setenv("BOOST_FACTORY_ECHO_PROTOCOL", "H2C")
	os.Setenv("BOOST_FACTORY_LOGRUS_CONSOLE_LEVEL", "TRACE")
}

func main() {

	boost.Start()

	ctx := context.Background()

	srv, _ := echo.NewServer(ctx,
		logplugin.Register,
		error_handler.Register)

	srv.GET("/", func(c e.Context) error {
		req := c.Request()
		format := `
			<code>
			  Protocol: %s<br>
			  Host: %s<br>
			  Remote Address: %s<br>
			  Method: %s<br>
			  Path: %s<br>
			</code>
		  `
		return c.HTML(http.StatusOK, fmt.Sprintf(format, req.Proto, req.Host, req.RemoteAddr, req.Method, req.URL.Path))
	})

	srv.Serve(ctx)

	// curl -v --http2-prior-knowledge http://localhost:8080
}
