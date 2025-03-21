package ollama

import (
	"context"
	"github.com/xgodev/boost/model/errors"
	"net/http"
	"net/url"

	"github.com/ollama/ollama/api"
	"github.com/xgodev/boost/wrapper/log"
)

// NewConnWithOptions registers a nats connection.
func NewConnWithOptions(ctx context.Context, options *Options) (*api.Client, error) {

	logger := log.FromContext(ctx)

	// Configurar o cliente Ollama
	baseURL, err := url.Parse(options.Endpoint)
	if err != nil {
		return nil, err
	}
	client := api.NewClient(baseURL, http.DefaultClient)
	if client == nil {
		return nil, errors.New("failed to create new client")
	}

	logger.Infof("Connected to Ollama server: %s", options.Endpoint)

	return client, nil
}

// NewConnWithConfigPath returns a new nats connection with options from config path.
func NewConnWithConfigPath(ctx context.Context, path string) (*api.Client, error) {
	options, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewConnWithOptions(ctx, options)
}

// NewConn returns a new connection with default options.
func NewConn(ctx context.Context) (*api.Client, error) {

	logger := log.FromContext(ctx)

	o, err := NewOptions()
	if err != nil {
		logger.Fatalf(err.Error())
	}

	return NewConnWithOptions(ctx, o)
}
