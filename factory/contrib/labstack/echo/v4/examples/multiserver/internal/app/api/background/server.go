// Package background is a server of app api that is not echo: any type with
// Serve and Shutdown runs next to the doors of the app.
package background

import (
	"context"
	"time"

	"github.com/xgodev/boost/extra/multiserver"
	"github.com/xgodev/boost/wrapper/log"
)

// Server serves for a while and stops, so GET /check reports it down.
type Server struct {
	o    *Options
	stop chan struct{}
}

// New provides the server twice: as *Server, which the rest plugins receive,
// and in the group of the servers that multiserver runs.
// @Provide(index=0)
// @Provide(group=extra.multiserver.servers, index=1)
func New(o *Options) (*Server, multiserver.Server) {
	s := &Server{o: o, stop: make(chan struct{})}
	return s, s
}

// Serve returns after the configured duration, or earlier when the app
// stops.
func (s *Server) Serve(ctx context.Context) {
	select {
	case <-time.After(s.o.Duration):
		log.FromContext(ctx).Infof("background server finished after %s", s.o.Duration)
	case <-s.stop:
	case <-ctx.Done():
	}
}

// Shutdown makes Serve return.
func (s *Server) Shutdown(context.Context) {
	close(s.stop)
}
