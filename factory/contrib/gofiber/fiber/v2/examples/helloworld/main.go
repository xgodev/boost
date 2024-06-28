package main

import (
	"context"
	"github.com/xgodev/boost/wrapper/config"
	"net/http"

	f "github.com/gofiber/fiber/v2"
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2/plugins/extra/error_handler"
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2/plugins/native/cors"
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2/plugins/native/etag"
	"github.com/xgodev/boost/factory/local/wrapper/log"
)

func Get(c *f.Ctx) (err error) {
	return c.Status(http.StatusOK).SendString("Hello!!")
}

func main() {

	config.Load()
	log.New()

	ctx := context.Background()

	srv := fiber.NewServer(ctx,
		error_handler.Register,
		cors.Register,
		etag.Register)

	srv.Get("/hello-world", Get)
	srv.Serve(ctx)
}
