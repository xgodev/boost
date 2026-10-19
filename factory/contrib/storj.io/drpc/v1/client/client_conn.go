package client

import (
	"context"
	"net"
	"strconv"

	"github.com/xgodev/boost/factory"
	"github.com/xgodev/boost/wrapper/log"
	"storj.io/drpc/drpcconn"
)

// Plugin defines a drpc client Plugin function to execute.
type Plugin = factory.Plugin[*drpcconn.Conn]

// NewClientConn returns a new drpc client connection.
func NewClientConn(ctx context.Context, plugins ...Plugin) (*drpcconn.Conn, error) {
	opt, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewClientConnWithOptions(ctx, opt, plugins...)
}

// NewClientConnWithConfigPath returns a new drpc client connection with options from config path.
func NewClientConnWithConfigPath(ctx context.Context, path string, plugins ...Plugin) (*drpcconn.Conn, error) {
	options, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewClientConnWithOptions(ctx, options, plugins...)
}

// NewClientConnWithOptions dials the server and returns a drpc client
// connection over it. The caller closes it. A plugin that fails closes the
// connection and returns the error.
func NewClientConnWithOptions(ctx context.Context, options *Options, plugins ...Plugin) (*drpcconn.Conn, error) {
	logger := log.FromContext(ctx)

	addr := net.JoinHostPort(options.Host, strconv.Itoa(options.Port))

	raw, err := (&net.Dialer{Timeout: options.Timeout}).DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, err
	}

	conn := drpcconn.New(raw)
	if err := factory.ApplyPlugins(ctx, conn, plugins...); err != nil {
		_ = conn.Close()
		return nil, err
	}

	logger.Debugf("drpc client created for host %s", addr)

	return conn, nil
}
