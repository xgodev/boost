package server

import (
	"github.com/xgodev/boost"
	"golang.org/x/net/http2"
)

// NewServer returns a pointer with new Server
func NewServer() (*http2.Server, error) {
	return boost.NewOptionsWithPath[http2.Server](root)
}

// NewServerWithPath returns a pointer with new Server
func NewServerWithPath(path string) (srv *http2.Server, err error) {
	return boost.NewOptionsWithPath[http2.Server](root, path)
}
