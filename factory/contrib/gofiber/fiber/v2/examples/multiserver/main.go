package main

import (
	"context"
	"fmt"
	"github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2"
	"net/http"
	"time"

	f "github.com/gofiber/fiber/v2"
	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/extra/multiserver"
	mserver "github.com/xgodev/boost/factory/contrib/gofiber/fiber/v2/plugins/local/multi-server"
	ilog "github.com/xgodev/boost/factory/local/wrapper/log"
)

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

	config.Load()

	ilog.New()

	ctx1 := context.Background()

	srv1 := fiber.NewServer(ctx1,
		mserver.Register)

	srv1.Get("/test", Get)

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
