package firestore

import (
	apiv1 "github.com/xgodev/boost/factory/contrib/cloud.google.com/api/v0"
	grpcv1 "github.com/xgodev/boost/factory/contrib/cloud.google.com/grpc/v1"
	"github.com/xgodev/boost/wrapper/config"
)

// Options holds shared API/gRPC options for Firestore.
type Options struct {
	APIOptions  apiv1.Options  `config:"apiOptions"`
	GRPCOptions grpcv1.Options `config:"grpcOptions"`
	Database    string         `config:"database"`
}

// NewOptions loads Options from the default root.
func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath loads Options from the specified path.
func NewOptionsWithPath(path string) (*Options, error) {
	return config.NewOptionsWithPath[Options](path)
}
