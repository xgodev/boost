package bigtable

import "github.com/xgodev/boost/wrapper/config"

const root = "boost.factory.bigtable"

func init() {
	ConfigAdd(root)
}

// ConfigAdd registers the Bigtable client keys under path, so one app can
// hold more than one client (boost.factory.bigtable.clients.<name>).
func ConfigAdd(path string) {
	config.Add(path+".project", "", "GCP project of the Bigtable instance")
	config.Add(path+".instance", "", "Bigtable instance ID")
	config.Add(path+".credentials.file", "", "service account key file; empty uses Application Default Credentials")
}
