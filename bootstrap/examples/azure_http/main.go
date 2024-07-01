package main

import (
	"context"
	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	logger "github.com/xgodev/boost/bootstrap/cloudevents/plugins/local/wrapper/log"
	"github.com/xgodev/boost/bootstrap/cmd"
	"github.com/xgodev/boost/config"
	igce "github.com/xgodev/boost/factory/contrib/cloudevents/sdk-go/v2"
	ilog "github.com/xgodev/boost/factory/local/wrapper/log"
	"go.uber.org/fx"
	"os"
)

func init() {
	// sets env var
	os.Setenv("FAAS_CMD_DEFAULT", "azure")
	os.Setenv("IGNITE_LOGRUS_CONSOLE_LEVEL", "TRACE")
}

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

	// go run main.go help
	err := cmd.Run(options,
		cmd.NewAzure(),
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

	randData := Message{
		Random: uuid.NewString(),
	}

	e.SetExtension("partitionkey", "changeme")

	err := e.SetData("", randData)

	return &e, err
}

type Message struct {
	Random string
}
