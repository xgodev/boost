package azure

import "github.com/xgodev/boost/wrapper/config"

const (
	// Root is the config root of the azure credential; the azure clients
	// live under it.
	Root         = "boost.factory.azure"
	clientId     = ".clientId"
	clientSecret = ".clientSecret"
	tenantId     = ".tenantId"
)

func init() {
	ConfigAdd(Root)
}

// ConfigAdd registers the azure credential config keys under path.
func ConfigAdd(path string) {
	config.Add(path+clientId, "", "defines the azure client id", config.WithHide())
	config.Add(path+clientSecret, "", "defines the azure client secret", config.WithHide())
	config.Add(path+tenantId, "", "defines the azure tenant id")
}
