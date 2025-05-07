package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	apiv1 "github.com/xgodev/boost/factory/contrib/cloud.google.com/api/v0"
	grpcv1 "github.com/xgodev/boost/factory/contrib/cloud.google.com/grpc/v1"
	clientgrpc "github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/client"
	"github.com/xgodev/boost/wrapper/log"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewClient creates a Firestore client using default configuration.
func NewClient(ctx context.Context, plugins ...clientgrpc.Plugin) (*firestore.Client, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewClientWithOptions(ctx, o, plugins...)
}

// NewClientWithConfigPath creates a Firestore client using configuration from the specified path.
func NewClientWithConfigPath(ctx context.Context, path string, plugins ...clientgrpc.Plugin) (*firestore.Client, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewClientWithOptions(ctx, o, plugins...)
}

// NewClientWithOptions constructs a Firestore client from Options.
func NewClientWithOptions(
	ctx context.Context,
	o *Options,
	plugins ...clientgrpc.Plugin,
) (*firestore.Client, error) {
	logger := log.FromContext(ctx)

	var clientOpts []option.ClientOption

	if o.APIOptions.UseEmulator {
		logger.Infof("using emulator at %s", o.APIOptions.EmulatorHost)
		clientOpts = []option.ClientOption{
			option.WithEndpoint(o.APIOptions.EmulatorHost),
			option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
			option.WithoutAuthentication(),
		}
	} else {
		apiOpts := apiv1.ApplyAPIOptions(ctx, &o.APIOptions)
		grpcDialOpts := grpcv1.ApplyDialOptions(ctx, &o.GRPCOptions, plugins...)

		clientOpts = make([]option.ClientOption, len(apiOpts))
		copy(clientOpts, apiOpts)
		for _, dop := range grpcDialOpts {
			clientOpts = append(clientOpts, option.WithGRPCDialOption(dop))
		}
	}

	logger.Debugf("creating Firestore client for project %s", o.APIOptions.ProjectID)
	return firestore.NewClient(ctx, o.APIOptions.ProjectID, clientOpts...)
}
