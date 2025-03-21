package ollama

import (
	"context"
	"github.com/xgodev/boost/model/errors"
	"net/http"
	"net/url"

	"github.com/ollama/ollama/api"
	"github.com/xgodev/boost/wrapper/log"
)

// NewClientWithOptions registers a ollama connection.
func NewClientWithOptions(ctx context.Context, options *Options) (*api.Client, error) {

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

// NewClientWithConfigPath returns a new nats connection with options from config path.
func NewClientWithConfigPath(ctx context.Context, path string) (*api.Client, error) {
	options, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewClientWithOptions(ctx, options)
}

// NewClient returns a new connection with default options.
func NewClient(ctx context.Context) (*api.Client, error) {

	logger := log.FromContext(ctx)

	o, err := NewOptions()
	if err != nil {
		logger.Fatalf(err.Error())
	}

	return NewClientWithOptions(ctx, o)
}
