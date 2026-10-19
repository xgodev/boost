package azblob

import (
	azure "github.com/xgodev/boost/factory/contrib/azure/azure-sdk-for-go/v1"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root             = azure.Root + ".azblob"
	accountName      = ".accountName"
	url              = ".url"
	connectionString = ".connectionString"
	// PluginsRoot is the config root every azblob plugin key lives under.
	PluginsRoot = root + ".plugins"
)

func init() {
	ConfigAdd(root)
}

// ConfigAdd registers the azblob client config keys under path.
func ConfigAdd(path string) {
	config.Add(path+accountName, "unknown", "defines the azure azblob account name", config.WithHide())
	config.Add(path+url, "", "defines the service url; empty means https://<accountName>.blob.core.windows.net/")
	config.Add(path+connectionString, "", "defines a connection string, used instead of the credential when set", config.WithHide())
}
