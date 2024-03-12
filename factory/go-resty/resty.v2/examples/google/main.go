package main

import (
	"context"

	r "github.com/go-resty/resty/v2"
	"github.com/xgodev/boost/config"
	"github.com/xgodev/boost/factory/go-resty/resty.v2"
	"github.com/xgodev/boost/factory/go-resty/resty.v2/plugins/contrib/americanas-go/health.v1"
	ilog "github.com/xgodev/boost/factory/xgodev/boost.v1/log"
	"github.com/xgodev/boost/log"
)

func main() {

	var err error

	config.Load()

	ctx := context.Background()

	ilog.New()

	logger := log.FromContext(ctx)

	options := health.Options{
		Name:        "Google Inc",
		Host:        "http://google.com",
		Endpoint:    "/status",
		Enabled:     true,
		Description: "Search Engine",
		Required:    true,
	}

	healthIntegrator := health.NewHealthWithOptions(&options)

	client := resty.NewClientWithOptions(ctx, &resty.Options{}, healthIntegrator.Register)
	request := client.R().EnableTrace()

	var resp *r.Response
	resp, err = request.Get("http://google.com")
	if err != nil {
		logger.Fatalf(err.Error())
	}

	if resp != nil {
		logger.Infof(resp.String())
	}
}
