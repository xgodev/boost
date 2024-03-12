package main

import (
	"context"
	"github.com/xgodev/boost"
	"os"

	"github.com/xgodev/boost/factory/contrib/rs/zerolog/v1"
	"github.com/xgodev/boost/factory/local/wrapper"
	"github.com/xgodev/boost/factory/local/wrapper/plugins/contrib/afex/hystrix-go/v0"
	logger "github.com/xgodev/boost/factory/local/wrapper/plugins/local/log"
	"github.com/xgodev/boost/log"
)

func init() {
	os.Setenv("IGNITE_ZEROLOG_LEVEL", "TRACE")
}

func main() {

	ctx := context.Background()

	boost.Boot()
	zerolog.NewLogger()

	var r string
	var err error

	wrp, _ := wrapper.NewAnyErrorWrapper[string](ctx, logger.NewAnyError[string], hystrix.NewAnyError[string])

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
