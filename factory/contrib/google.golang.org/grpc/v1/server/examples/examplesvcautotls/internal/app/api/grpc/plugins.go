// Package grpc holds what the grpc server of app api needs besides its
// services.
package grpc

import (
	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/server"
	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/server/plugins/local/wrapper/log"
)

// Plugins log the calls.
// @Provide(index=0)
func Plugins() []server.Plugin {
	return []server.Plugin{log.Register}
}
