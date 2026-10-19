package client

import (
	"context"

	"github.com/valyala/fasthttp"
	"github.com/xgodev/boost/factory"
	"github.com/xgodev/boost/wrapper/log"
)

// Plugin defines a fasthttp client Plugin function to execute.
type Plugin = factory.Plugin[*fasthttp.Client]

// NewClient returns a new fasthttp client with default options.
func NewClient(ctx context.Context, plugins ...Plugin) (*fasthttp.Client, error) {
	opts, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewClientWithOptions(ctx, opts, plugins...)
}

// NewClientWithConfigPath returns a new fasthttp client with options from config path.
func NewClientWithConfigPath(ctx context.Context, path string, plugins ...Plugin) (*fasthttp.Client, error) {
	opts, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewClientWithOptions(ctx, opts, plugins...)
}

// NewClientWithOptions returns a new fasthttp client with options. A plugin
// that fails stops the build: the error is returned and later plugins do not
// run.
func NewClientWithOptions(ctx context.Context, o *Options, plugins ...Plugin) (*fasthttp.Client, error) {
	logger := log.FromContext(ctx)

	client := &fasthttp.Client{
		Name:                          o.Name,
		NoDefaultUserAgentHeader:      o.NoDefaultUserAgentHeader,
		DialDualStack:                 o.DialDualStack,
		MaxConnsPerHost:               o.MaxConnsPerHost,
		MaxIdleConnDuration:           o.MaxIdleConnDuration,
		MaxConnDuration:               o.MaxConnDuration,
		MaxIdemponentCallAttempts:     o.MaxIdemponentCallAttempts,
		ReadBufferSize:                o.ReadBufferSize,
		WriteBufferSize:               o.WriteBufferSize,
		ReadTimeout:                   o.ReadTimeout,
		WriteTimeout:                  o.WriteTimeout,
		MaxResponseBodySize:           o.MaxResponseBodySize,
		DisableHeaderNamesNormalizing: o.DisableHeaderNamesNormalizing,
		MaxConnWaitTimeout:            o.MaxConnWaitTimeout,
	}

	if err := factory.ApplyPlugins(ctx, client, plugins...); err != nil {
		return nil, err
	}

	logger.Debugf("fasthttp client %s created", o.Name)

	return client, nil
}
