package multiserver

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"

	server "github.com/xgodev/boost/extra/multiserver"
	contextfx "github.com/xgodev/boost/fx/modules/core/context"
	"github.com/xgodev/boost/wrapper/log"
	"go.uber.org/fx"
)

const (
	ServersGroupKey = "extra.multiserver.servers"
)

type params struct {
	fx.In
	Lifecycle  fx.Lifecycle
	Shutdowner fx.Shutdowner
	Servers    []server.Server `group:"extra.multiserver.servers"`
}

var once sync.Once

// Module fx module for multiserver. The servers start on the app's OnStart
// and drain on its OnStop, so app.Run handles SIGTERM gracefully.
func Module() fx.Option {
	options := fx.Options()

	once.Do(func() {

		options = fx.Options(
			contextfx.Module(),
			fx.Invoke(register),
		)
	})

	return options
}

func register(ctx context.Context, p params) error {
	if len(p.Servers) == 0 {
		return errors.New("multiserver: no servers configured")
	}

	serveCtx, cancel := context.WithCancel(ctx)
	served := make(chan struct{})
	var stopping atomic.Bool

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				defer close(served)
				server.Serve(serveCtx, p.Servers...)
				if stopping.Load() {
					return
				}
				logger := log.FromContext(ctx)
				logger.Error("multiserver: every server stopped unexpectedly, shutting down")
				if err := p.Shutdowner.Shutdown(fx.ExitCode(1)); err != nil {
					logger.Error(err)
				}
			}()
			return nil
		},
		OnStop: func(stopCtx context.Context) error {
			stopping.Store(true)
			defer cancel()

			drained := make(chan struct{})
			go func() {
				var wg sync.WaitGroup
				for _, srv := range p.Servers {
					wg.Go(func() { srv.Shutdown(stopCtx) })
				}
				wg.Wait()
				<-served
				close(drained)
			}()

			select {
			case <-drained:
				return nil
			case <-stopCtx.Done():
				return stopCtx.Err()
			}
		},
	})

	return nil
}
