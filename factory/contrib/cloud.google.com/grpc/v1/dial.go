package grpc

import (
	"context"
	"crypto/tls"
	"github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// ApplyDialOptions retorna os DialOptions gRPC baseados em Options
// e qualquer plugin passado como parâmetro.
func ApplyDialOptions(ctx context.Context, o *Options, plugins ...client.Plugin) []grpc.DialOption {
	var opts []grpc.DialOption
	// TLS or insecure
	if o.TLS.Enabled {
		// load TLS creds
		creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: o.TLS.InsecureSkipVerify})
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// window sizes
	/*
		opts = append(opts,
			grpc.WithInitialWindowSize(int32(o.InitialWindowSize)),
			grpc.WithInitialConnWindowSize(int32(o.InitialConnWindowSize)),
		)
	*/

	// authority override
	/*
		if o.HostOverwrite != "" {
			opts = append(opts, grpc.WithAuthority(o.HostOverwrite))
		}
	*/

	// backoff
	/*
		opts = append(opts, grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  o.Backoff.BaseDelay,
				Multiplier: o.Backoff.Multiplier,
				Jitter:     o.Backoff.Jitter,
				MaxDelay:   o.Backoff.MaxDelay,
			},
			MinConnectTimeout: o.MinConnectTimeout,
		}))
	*/

	// keepalive
	/*
		opts = append(opts, grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                o.Keepalive.Time,
			Timeout:             o.Keepalive.Timeout,
			PermitWithoutStream: o.Keepalive.PermitWithoutStream,
		}))
	*/

	// plugins de gRPC
	for _, plugin := range plugins {
		dopts, _ := plugin(ctx)
		opts = append(opts, dopts...)
	}

	return opts
}
