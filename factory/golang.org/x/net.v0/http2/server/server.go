package server

import (
	"github.com/xgodev/boost/factory"
	"golang.org/x/net/http2"
)

// NewServer returns a pointer with new Server
func NewServer() (*http2.Server, error) {
	return factory.NewOptionsWithPath[http2.Server](root)
}

// NewServerWithPath returns a pointer with new Server
func NewServerWithPath(path string) (srv *http2.Server, err error) {
	return factory.NewOptionsWithPath[http2.Server](root, path)
}
