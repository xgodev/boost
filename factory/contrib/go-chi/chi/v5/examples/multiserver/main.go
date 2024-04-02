package main

import (
	"context"
	"fmt"
	"github.com/xgodev/boost"
	"net/http"
	"time"

	"github.com/xgodev/boost/extra/multiserver"
	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5"
	multiserverplugin "github.com/xgodev/boost/factory/contrib/go-chi/chi/v5/plugins/local/extra/multiserver"
	"github.com/xgodev/boost/factory/contrib/go-chi/chi/v5/plugins/native/recoverer"
	"github.com/xgodev/boost/factory/local/wrapper/log"
)

func Get(ctx context.Context) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func main() {

	boost.Start()
	log.New()

	chiCtx := context.Background()
	chiSrv := chi.NewServer(chiCtx,
		multiserverplugin.Register,
		recoverer.Register,
	)

	chiSrv.Get("/hello", Get(chiCtx))

	msCtx := context.Background()
	multiserver.Serve(msCtx, chiSrv, &LocalServer{})
}

type LocalServer struct {
}

func (s *LocalServer) Serve(ctx context.Context) {
	time.Sleep(30 * time.Second)
	fmt.Printf("finished")
}

func (s *LocalServer) Shutdown(ctx context.Context) {
}
