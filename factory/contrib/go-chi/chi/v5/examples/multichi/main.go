package main

import (
	"context"
	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5"
	"github.com/xgodev/boost/wrapper/config"
	"net/http"

	"github.com/xgodev/boost/extra/multiserver"
	multiserverplugin "github.com/xgodev/boost/factory/contrib/go-chi/chi/v5/plugins/local/extra/multiserver"
	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5/plugins/native/recoverer"
	"github.com/xgodev/boost/factory/core/net/http/server"
	"github.com/xgodev/boost/factory/local/wrapper/log"
)

const httpServerRoot = "boost.factory.http2.server"

func init() {
	server.ConfigAdd(httpServerRoot)
}

func Get(ctx context.Context) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func main() {

	config.Load()
	log.New()

	chi1Ctx := context.Background()
	chi1Srv := chi.NewServer(chi1Ctx,
		multiserverplugin.Register,
		recoverer.Register,
	)

	chi1Srv.Get("/hello", Get(chi1Ctx))

	srv2Options, err := server.NewOptionsWithPath(httpServerRoot)
	if err != nil {
		panic(err)
	}

	chi2Ctx := context.Background()
	chi2Srtv := chi.NewServerWithOptions(chi2Ctx, srv2Options)

	msCtx := context.Background()
	multiserver.Serve(msCtx, chi1Srv, chi2Srtv)
}
