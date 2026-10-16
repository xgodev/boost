package multiserver

import (
	"context"
	"sync"

	"github.com/xgodev/boost/model/errors"
)

type d struct {
	Ok  bool
	Srv Server
}

// s records every served server; Check reads it from request goroutines
// while Serve's goroutines mark servers down, so mu guards it.
var (
	mu sync.Mutex
	s  []d
)

type Server interface {
	Serve(context.Context)
	Shutdown(context.Context)
}

func Serve(ctx context.Context, srvs ...Server) {

	switch len(srvs) {
	case 0:
		panic("no servers configured")
	case 1:
		i := add(srvs[0])
		srvs[0].Serve(ctx)
		markDown(i)
	default:
		wg := new(sync.WaitGroup)

		for _, srv := range srvs {
			i := add(srv)
			wg.Go(func() {
				srv.Serve(ctx)
				markDown(i)
			})
		}

		wg.Wait()
	}

}

func add(srv Server) int {
	mu.Lock()
	defer mu.Unlock()
	s = append(s, d{
		Ok:  true,
		Srv: srv,
	})
	return len(s) - 1
}

func markDown(i int) {
	mu.Lock()
	defer mu.Unlock()
	s[i].Ok = false
}

func snapshot() []d {
	mu.Lock()
	defer mu.Unlock()
	if len(s) == 0 {
		panic("no servers configured")
	}
	return append([]d(nil), s...)
}

func Check(ctx context.Context) error {
	for _, a := range snapshot() {
		if !a.Ok {
			return errors.ServiceUnavailablef("one of servers is down")
		}
	}

	return nil
}

func Shutdown(ctx context.Context) {
	for _, a := range snapshot() {
		a.Srv.Shutdown(ctx)
	}
}
