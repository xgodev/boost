package main

import (
	"context"
	"github.com/xgodev/boost"
	"net/http"

	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5"
	datadog "github.com/xgodev/boost/factory/contrib/go-chi/chi/v5/plugins/contrib/datadog/dd-trace-go/v1"
	newrelic "github.com/xgodev/boost/factory/contrib/go-chi/chi/v5/plugins/contrib/newrelic/go-agent/v3"
	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5/plugins/local/extra/health"
	multiserverplugin "github.com/xgodev/boost/factory/contrib/go-chi/chi/v5/plugins/local/extra/multiserver"
	status "github.com/xgodev/boost/factory/contrib/go-chi/chi/v5/plugins/local/model/restresponse"
	logplugin "github.com/xgodev/boost/factory/contrib/go-chi/chi/v5/plugins/local/wrapper/log"
	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5/plugins/native/realip"
	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5/plugins/native/recoverer"
	ifx "github.com/xgodev/boost/factory/contrib/go.uber.org/fx/v1"
	fxchi "github.com/xgodev/boost/factory/contrib/go.uber.org/fx/v1/module/contrib/go-chi/chi/v5"
	fxctx "github.com/xgodev/boost/factory/contrib/go.uber.org/fx/v1/module/core/context"
	ilog "github.com/xgodev/boost/factory/local/wrapper/log"
	"github.com/xgodev/boost/wrapper/log"
	"go.uber.org/fx"
)

func Get(ctx context.Context) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func main() {

	boost.Start()
	ilog.New()

	ifx.NewApp(
		fxctx.Module(),
		fxchi.Module(),
		fx.Provide(
			func() []chi.Plugin {
				return []chi.Plugin{
					multiserverplugin.Register,
					recoverer.Register,
					realip.Register,
					logplugin.Register,
					status.Register,
					health.Register,
					newrelic.Register,
					datadog.Register,
				}
			},
		),
		fx.Invoke(
			func(server *chi.Server, ctx context.Context) {
				server.Get("/hello", Get(ctx))
			},
		),
		fx.Invoke(
			func(lifecycle fx.Lifecycle, server *chi.Server) {
				lifecycle.Append(
					fx.Hook{
						OnStart: func(ctx context.Context) error {
							log.Info("starting server")
							go server.Serve(ctx)
							return nil
						},
						OnStop: func(ctx context.Context) error {
							log.Info("stopping server")
							server.Shutdown(ctx)
							return nil
						},
					},
				)
			},
		),
	).Run()

}
