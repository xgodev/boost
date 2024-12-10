package pubsub

import (
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root            = "boost.factory.pubsub"
	projectID       = ".projectId"
	credentialsRoot = ".credentials"
	credentialsFile = credentialsRoot + ".file"
	credentialsJSON = credentialsRoot + ".json"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+projectID, "default", "defines project ID")
	config.Add(path+credentialsFile, "", "sets credentials file")
	config.Add(path+credentialsJSON, "", "sets credentials json")
}
