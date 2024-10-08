package resty

import (
	"context"
	"net"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/xgodev/boost/wrapper/log"
)

// Plugin defines a function to process plugin.
type Plugin func(context.Context, *resty.Client) error

// NewClient returns a new resty Client.
func NewClient(ctx context.Context, plugins ...Plugin) (*resty.Client, error) {
	opts, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewClientWithOptions(ctx, opts, plugins...), nil
}

// NewClientWithConfigPath returns a new resty Client with options from config path.
func NewClientWithConfigPath(ctx context.Context, path string, plugins ...Plugin) (*resty.Client, error) {
	opts, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewClientWithOptions(ctx, opts, plugins...), nil
}

// NewClientWithOptions returns a new resty Client with options.
func NewClientWithOptions(ctx context.Context, options *Options, plugins ...Plugin) *resty.Client {

	logger := log.FromContext(ctx)

	logger.Tracef("creating resty client for host %s", options.Host)

	client := resty.New()

	dialer := &net.Dialer{
		Timeout:       options.ConnectionTimeout,
		FallbackDelay: options.FallbackDelay,
		KeepAlive:     options.KeepAlive,
	}

	transport := &http.Transport{
		DisableCompression:    options.Transport.DisableCompression,
		DisableKeepAlives:     options.Transport.DisableKeepAlives,
		MaxIdleConnsPerHost:   options.Transport.MaxIdleConnsPerHost,
		ResponseHeaderTimeout: options.Transport.ResponseHeaderTimeout,
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     options.Transport.ForceAttemptHTTP2,
		MaxIdleConns:          options.Transport.MaxIdleConns,
		MaxConnsPerHost:       options.Transport.MaxConnsPerHost,
		IdleConnTimeout:       options.Transport.IdleConnTimeout,
		TLSHandshakeTimeout:   options.Transport.TLSHandshakeTimeout,
		ExpectContinueTimeout: options.Transport.ExpectContinueTimeout,
	}

	client.
		SetTransport(transport).
		SetHeader("accept", options.Accept).
		SetTimeout(options.RequestTimeout).
		SetDebug(options.Debug).
		SetBaseURL(options.Host).
		SetCloseConnection(options.CloseConnection)

	if options.Authorization != "" {
		client.SetHeader("Authorization", options.Authorization)
	}

	for k, v := range options.Headers {
		client.SetHeader(k, v)
	}

	for k, v := range options.QueryParams {
		client.SetQueryParam(k, v)
	}

	for _, plugin := range plugins {
		if err := plugin(ctx, client); err != nil {
			panic(err)
		}
	}

	logger.Debugf("resty client created for host %s", options.Host)

	return client
}
