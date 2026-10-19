package sns

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	awsfactory "github.com/xgodev/boost/factory/contrib/aws/aws-sdk-go-v2/v1"
)

// ConfigAdd adds the aws keys under path, for an SDK client with a config
// of its own.
func ConfigAdd(path string) {
	awsfactory.ConfigAdd(path)
}

// NewSDKClient returns the SNS client of the SDK, built from the aws keys
// under boost.factory.aws.
func NewSDKClient(ctx context.Context, plugins ...awsfactory.Plugin) (*sns.Client, error) {
	cfg, err := awsfactory.NewConfig(ctx, plugins...)
	if err != nil {
		return nil, err
	}
	return sns.NewFromConfig(cfg), nil
}

// NewSDKClientWithConfigPath returns the SNS client of the SDK, built from
// the aws keys under path.
func NewSDKClientWithConfigPath(ctx context.Context, path string, plugins ...awsfactory.Plugin) (*sns.Client, error) {
	cfg, err := awsfactory.NewConfigWithConfigPath(ctx, path, plugins...)
	if err != nil {
		return nil, err
	}
	return sns.NewFromConfig(cfg), nil
}

// NewSDKClientWithOptions returns the SNS client of the SDK, built from
// options.
func NewSDKClientWithOptions(ctx context.Context, options *awsfactory.Options, plugins ...awsfactory.Plugin) (*sns.Client, error) {
	cfg, err := awsfactory.NewConfigWithOptions(ctx, options, plugins...)
	if err != nil {
		return nil, err
	}
	return sns.NewFromConfig(cfg), nil
}
