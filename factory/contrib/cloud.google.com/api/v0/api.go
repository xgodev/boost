package api

import (
	"context"
	"github.com/xgodev/boost/wrapper/log"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"net/url"
	"os"
)

// ApplyAPIOptions retorna os option.ClientOption da biblioteca
// google.golang.org/api baseados em Options.
func ApplyAPIOptions(ctx context.Context, o *Options) []option.ClientOption {

	logger := log.FromContext(ctx)

	var opts []option.ClientOption
	// proxy
	if o.Proxy != "" {
		if u, err := url.Parse(o.Proxy); err == nil {
			httpc := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(u)}}
			opts = append(opts, option.WithHTTPClient(httpc))
		}
	}
	// emulator
	if o.UseEmulator {
		logger.Debugf("using emulator at %s", o.EmulatorHost)
		host := o.EmulatorHost
		if host == "" {
			host = os.Getenv("EMULATOR_HOST")
		}
		opts = append(opts, option.WithEndpoint(host), option.WithoutAuthentication(), option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())))

	} else {
		// credentials
		if o.Credentials.JSON != "" {
			opts = append(opts, option.WithCredentialsJSON([]byte(o.Credentials.JSON)))
		} else if o.Credentials.File != "" {
			opts = append(opts, option.WithCredentialsFile(o.Credentials.File))
		}
	}
	// endpoint, scopes, user-agent
	if o.Endpoint != "" {
		opts = append(opts, option.WithEndpoint(o.Endpoint))
	}
	if len(o.ParseScopes()) > 0 {
		opts = append(opts, option.WithScopes(o.ParseScopes()...))
	}
	if o.UserAgent != "" {
		opts = append(opts, option.WithUserAgent(o.UserAgent))
	}

	// opts = append(opts, option.WithLogger(log.GetLogger()))

	return opts
}
