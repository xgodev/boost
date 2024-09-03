package main

import (
	"context"
	"github.com/xgodev/boost"
	"os"

	"github.com/xgodev/boost/extra/middleware"
	"github.com/xgodev/boost/extra/middleware/plugins/contrib/afex/hystrix-go/v0"
	logger "github.com/xgodev/boost/extra/middleware/plugins/local/wrapper/log"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
)

func init() {
	os.Setenv("BOOST_FACTORY_ZEROLOG_LEVEL", "TRACE")
}

func main() {

	ctx := context.Background()

	boost.Start()

	var r string
	var err error

	wrp := middleware.NewAnyErrorWrapper[string](ctx, "test", logger.NewAnyErrorMiddleware[string](ctx), hystrix.NewAnyErrorMiddleware[string](ctx, "test"))

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
