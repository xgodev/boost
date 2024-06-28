package main

import (
	"context"
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
	"github.com/xgodev/boost/wrapper/config"
	"net/http"
	"os"

	f "github.com/gofiber/fiber/v2"
	"github.com/xgodev/boost/extra/multiserver"
	ilog "github.com/xgodev/boost/factory/local/wrapper/log"
)

func init() {
	fiber.ConfigAdd("boost.factory.fiber2")
}

type Response struct {
	Message string
}

func Get(c *f.Ctx) (err error) {

	resp := Response{
		Message: "Hello World!!",
	}

	err = config.Unmarshal(&resp)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(resp)
}

func main() {

	os.Setenv("BOOST_FACTORY_FIBER2_PORT", "8086")

	config.Load()

	ilog.New()

	ctx1 := context.Background()

	srv1 := fiber.NewServer(ctx1)

	srv1.Get("/test", Get)

	ctx2 := context.Background()

	srv2, err := fiber.NewServerWithConfigPath(ctx2, "boost.factory.fiber2")
	if err != nil {
		panic(err)
	}

	srv2.Get("/test", Get)

	multiserver.Serve(context.Background(), srv1, srv2)
}
