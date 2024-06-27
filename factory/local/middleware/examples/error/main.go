package main

import (
	"context"
	"github.com/xgodev/boost"
	"os"

	"github.com/xgodev/boost/factory/contrib/rs/zerolog/v1"
	"github.com/xgodev/boost/factory/local/middleware"
	"github.com/xgodev/boost/factory/local/middleware/plugins/contrib/afex/hystrix-go/v0"
	logger "github.com/xgodev/boost/factory/local/middleware/plugins/local/wrapper/log"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
)

func init() {
	os.Setenv("BOOST_FACTORY_ZEROLOG_LEVEL", "TRACE")
}

func main() {

	ctx := context.Background()

	boost.Start()
	zerolog.NewLogger()

	var r string
	var err error

	wrp, _ := middleware.NewAnyErrorWrapper[string](ctx, logger.NewAnyError[string], hystrix.NewAnyError[string])

	r, err = wrp.Exec(ctx, "xpto",
		func(ctx context.Context) (string, error) {
			l := log.FromContext(ctx)
			l.Info("executed business rule with error")
			return "", errors.New("an error ocurred")
		},
		func(ctx context.Context, v string, err error) (string, error) {
			l := log.FromContext(ctx)
			if err != nil {
				l.Info("executed fallback business rule")
				return "string", nil
			}
			return "", err
		})

	if err != nil {
		log.Errorf(err.Error())
	}

	log.Infof(r)
}
