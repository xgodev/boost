package pubsub

import (
	apiv1 "github.com/xgodev/boost/factory/contrib/cloud.google.com/api/v0"
	grpcv1 "github.com/xgodev/boost/factory/contrib/cloud.google.com/grpc/v1"
)

const root = "boost.factory.gcp.pubsub"

func init() {
	ConfigAdd(root)
}

// ConfigAdd registers shared API and gRPC configs under the given path.
func ConfigAdd(path string) {
	apiv1.ConfigAdd(path + ".apiOptions")
	grpcv1.ConfigAdd(path + ".grpcOptions")
}
