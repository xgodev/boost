package main

import (
	"context"
	r "github.com/go-resty/resty/v2"
	"github.com/xgodev/boost"
	"github.com/xgodev/boost/factory/contrib/go-resty/resty/v2"
	"github.com/xgodev/boost/factory/contrib/go-resty/resty/v2/plugins/local/extra/health"
	ilog "github.com/xgodev/boost/factory/local/wrapper/log"
	"github.com/xgodev/boost/wrapper/log"
)

func main() {

	var err error

	boost.Start()

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
