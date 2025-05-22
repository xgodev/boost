package bigquery

import (
	apiv1 "github.com/xgodev/boost/factory/contrib/cloud.google.com/api/v0"
	grpcv1 "github.com/xgodev/boost/factory/contrib/cloud.google.com/grpc/v1"
	"github.com/xgodev/boost/wrapper/config"
)

const root = "boost.factory.gcp.bigquery"

func init() {
	ConfigAdd(root)
}

// ConfigAdd registers BigQuery-specific and shared GCP configs under the given root.
func ConfigAdd(path string) {
	// shared API and gRPC configs
	apiv1.ConfigAdd(path + ".apiOptions")
	grpcv1.ConfigAdd(path + ".grpcOptions")

	// BigQuery-specific settings
	config.Add(path+".userProject", "", "alternative billing project (UserProject)")
}
