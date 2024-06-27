package main

import (
	"context"
	logger "github.com/xgodev/boost/bootstrap/cloudevents/plugins/local/wrapper/log"
	"os"

	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/xgodev/boost/bootstrap/cmd"
	"github.com/xgodev/boost/config"
	igce "github.com/xgodev/boost/factory/contrib/cloudevents/sdk-go/v2"
	ilog "github.com/xgodev/boost/factory/local/wrapper/log"
	"github.com/xgodev/boost/wrapper/log"
	"go.uber.org/fx"
)

func main() {

	config.Load()
	ilog.New()

	options := fx.Options(
		logger.Module(),
		fx.Provide(
			func() igce.Handler {
				return Handle
			},
		),
	)

	// sets env var
	os.Setenv("FAAS_CMD_DEFAULT", "cloudevents")

	// go run main.go help
	err := cmd.Run(options,
		cmd.NewCloudEvents(),
	)

	if err != nil {
		panic(err)
	}
}

func Handle(ctx context.Context, in v2.Event) (*v2.Event, error) {

	log.Info(in.Data())

	return nil, nil
}
