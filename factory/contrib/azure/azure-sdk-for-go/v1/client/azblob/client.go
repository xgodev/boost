package azblob

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/xgodev/boost/factory"
	"github.com/xgodev/boost/wrapper/log"
)

// Plugin customizes the azblob client options before the client is built.
type Plugin = factory.Plugin[*azblob.ClientOptions]

// NewClient returns an azblob client with default options.
func NewClient(ctx context.Context, credential azcore.TokenCredential, plugins ...Plugin) (*azblob.Client, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewClientWithOptions(ctx, credential, o, plugins...)
}

// NewClientWithConfigPath returns an azblob client with options from config path.
func NewClientWithConfigPath(ctx context.Context, credential azcore.TokenCredential, path string, plugins ...Plugin) (*azblob.Client, error) {
	opts, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewClientWithOptions(ctx, credential, opts, plugins...)
}

// NewClientWithOptions returns an azblob client with options. With a
// connection string it ignores credential; otherwise it authenticates with
// credential against the url, or the account's public url when the url is
// empty. A plugin that fails stops the build and its error is returned.
func NewClientWithOptions(ctx context.Context, credential azcore.TokenCredential, options *Options, plugins ...Plugin) (*azblob.Client, error) {
	logger := log.FromContext(ctx)

	clientOptions := &azblob.ClientOptions{}
	if err := factory.ApplyPlugins(ctx, clientOptions, plugins...); err != nil {
		return nil, err
	}

	if options.ConnectionString != "" {
		logger.Debugf("creating azblob client from connection string")
		return azblob.NewClientFromConnectionString(options.ConnectionString, clientOptions)
	}

	url := options.URL
	if url == "" {
		url = fmt.Sprintf("https://%s.blob.core.windows.net/", options.AccountName)
	}
	logger.Debugf("creating azblob client for %s", url)

	return azblob.NewClient(url, credential, clientOptions)
}
