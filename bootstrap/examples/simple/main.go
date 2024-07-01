package main

import (
	"context"
	"github.com/xgodev/boost/bootstrap/cloudevents/plugins/local/wrapper/log"
	"os"

	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/xgodev/boost/bootstrap/cloudevents/plugins/extra/publisher"
	"github.com/xgodev/boost/bootstrap/cmd"
	"github.com/xgodev/boost/config"
	igce "github.com/xgodev/boost/factory/contrib/cloudevents/sdk-go/v2"
	ilog "github.com/xgodev/boost/factory/local/wrapper/log"
	"go.uber.org/fx"
)

func main() {

	config.Load()

	ilog.New()

	options := fx.Options(
		log.Module(),
		publisher.Module(),
		fx.Provide(
			func() igce.Handler {
				return Handle
			},
		),
	)

	// sets env var
	os.Setenv("FAAS_CMD_DEFAULT", "nats")

	// go run main.go help
	err := cmd.Run(options,

		// go run main.go nats
		// or
		// FAAS_CMD_DEFAULT=nats go run main.go
		cmd.NewNats(),

		// go run main.go cloudevents
		// or
		// FAAS_CMD_DEFAULT=cloudevents go run main.go
		cmd.NewCloudEvents(),

		// go run main.go lambda
		// or
		// FAAS_CMD_DEFAULT=lambda go run main.go
		cmd.NewLambda(),
	)

	if err != nil {
		panic(err)
	}
}

func Handle(ctx context.Context, in v2.Event) (*v2.Event, error) {

	e := v2.NewEvent()
	e.SetID("changeme")
	e.SetSubject("changeme")
	e.SetSource("changeme")
	e.SetType("changeme")
	e.SetExtension("partitionkey", "changeme")
	err := e.SetData(v2.TextPlain, "changeme")

	return &e, err
}
