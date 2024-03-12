package main

import (
	"context"
	"os"

	"github.com/xgodev/boost/factory"
	"github.com/xgodev/boost/factory/rs/zerolog.v1"
	"github.com/xgodev/boost/factory/xgodev/boost.v1/grapper"
	"github.com/xgodev/boost/factory/xgodev/boost.v1/grapper/plugins/contrib/afex/hystrix-go.v0"
	logger "github.com/xgodev/boost/factory/xgodev/boost.v1/grapper/plugins/contrib/xgodev/boost.v1/log"
	"github.com/xgodev/boost/log"
)

func init() {
	os.Setenv("IGNITE_ZEROLOG_LEVEL", "TRACE")
}

func main() {

	ctx := context.Background()

	factory.Boot()
	zerolog.NewLogger()

	var r string
	var err error

	wrp, _ := grapper.NewAnyErrorWrapper[string](ctx, logger.NewAnyError[string], hystrix.NewAnyError[string])

	r, err = wrp.Exec(ctx, "xpto",
		func(ctx context.Context) (string, error) {
			l := log.FromContext(ctx)
			l.Info("executed business rule")
			return "string", nil
		}, nil)

	if err != nil {
		log.Errorf(err.Error())
	}

	log.Infof(r)
}
