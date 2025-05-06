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
func NewClientWithOptions(ctx context.Context, o *Options, plugins ...clientgrpc.Plugin) (*firestore.Client, error) {
	logger := log.FromContext(ctx)

	// aplica opções de API leve
	apiOpts := apiv1.ApplyAPIOptions(ctx, &o.APIOptions)
	// aplica opções de dial gRPC (plugins, interceptors etc)
	grpcDialOpts := grpcv1.ApplyDialOptions(ctx, &o.GRPCOptions, plugins...)

	// monta clientOpts
	clientOpts := make([]option.ClientOption, 0, len(apiOpts)+len(grpcDialOpts)+2)
	clientOpts = append(clientOpts, apiOpts...)
	for _, dop := range grpcDialOpts {
		clientOpts = append(clientOpts, option.WithGRPCDialOption(dop))
	}
	for _, plugin := range plugins {
		if dopts, err := plugin(ctx); err == nil {
			for _, dop := range dopts {
				clientOpts = append(clientOpts, option.WithGRPCDialOption(dop))
			}
		}
	}

	// se detectar o emulator, override do endpoint e desliga TLS/auth
	if o.APIOptions.EmulatorHost != "" && o.APIOptions.UseEmulator {
		logger.Infof("using Firestore emulator at %s", o.APIOptions.EmulatorHost)
		clientOpts = append(clientOpts,
			option.WithEndpoint(o.APIOptions.EmulatorHost), // host:porta do emulator
			option.WithGRPCDialOption(grpc.WithTransportCredentials( // força canal inseguro
				insecure.NewCredentials(),
			)),
			option.WithoutAuthentication(), // não tenta credenciais GCP
		)
	}

	logger.Debugf("creating Firestore client for project %s", o.APIOptions.ProjectID)
	return firestore.NewClient(ctx, o.APIOptions.ProjectID, clientOpts...)
}
