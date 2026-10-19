package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/xgodev/boost/wrapper/log"
)

// NewCredential returns an azure token credential with default options.
func NewCredential(ctx context.Context) (azcore.TokenCredential, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewCredentialWithOptions(ctx, o)
}

// NewCredentialWithConfigPath returns an azure token credential with options from config path.
func NewCredentialWithConfigPath(ctx context.Context, path string) (azcore.TokenCredential, error) {
	opts, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewCredentialWithOptions(ctx, opts)
}

// NewCredentialWithOptions returns a client secret credential when client
// id, client secret and tenant id are all set, and the default azure
// credential chain (environment, workload identity, managed identity, az
// cli) otherwise.
func NewCredentialWithOptions(ctx context.Context, options *Options) (azcore.TokenCredential, error) {
	logger := log.FromContext(ctx)

	if options.ClientId != "" && options.ClientSecret != "" && options.TenantId != "" {
		logger.Debugf("creating azure client secret credential for tenant %s", options.TenantId)
		return azidentity.NewClientSecretCredential(options.TenantId, options.ClientId, options.ClientSecret, nil)
	}

	logger.Debugf("creating azure default credential")
	return azidentity.NewDefaultAzureCredential(nil)
}
