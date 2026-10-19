package bigtable

import (
	"context"

	"cloud.google.com/go/bigtable"
	clientgrpc "github.com/xgodev/boost/factory/contrib/google.golang.org/grpc/v1/client"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
	"google.golang.org/api/option"
)

// NewClient creates a Bigtable client from boost.factory.bigtable.
func NewClient(ctx context.Context, plugins ...clientgrpc.Plugin) (*bigtable.Client, error) {
	o, err := NewOptions()
	if err != nil {
		return nil, err
	}
	return NewClientWithOptions(ctx, o, plugins...)
}

// NewClientWithConfigPath creates a Bigtable client from the keys under path.
func NewClientWithConfigPath(ctx context.Context, path string, plugins ...clientgrpc.Plugin) (*bigtable.Client, error) {
	o, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, err
	}
	return NewClientWithOptions(ctx, o, plugins...)
}

// NewClientWithOptions creates a Bigtable client from Options. Project and
// instance are required: without them the SDK would fail only on its first
// call to Google. Each plugin adds gRPC dial options. With
// BIGTABLE_EMULATOR_HOST set, the SDK dials the emulator itself, without
// credentials and without the plugins' dial options.
func NewClientWithOptions(ctx context.Context, o *Options, plugins ...clientgrpc.Plugin) (*bigtable.Client, error) {
	if o.Project == "" || o.Instance == "" {
		return nil, errors.NotValidf("bigtable: project %q and instance %q are both required", o.Project, o.Instance)
	}

	var opts []option.ClientOption
	if o.Credentials.File != "" {
		opts = append(opts, option.WithAuthCredentialsFile(option.ServiceAccount, o.Credentials.File))
	}
	for _, plugin := range plugins {
		dialOpts, _ := plugin(ctx)
		for _, dop := range dialOpts {
			opts = append(opts, option.WithGRPCDialOption(dop))
		}
	}

	log.FromContext(ctx).Debugf("creating Bigtable client for project %s, instance %s", o.Project, o.Instance)
	return bigtable.NewClient(ctx, o.Project, o.Instance, opts...)
}
