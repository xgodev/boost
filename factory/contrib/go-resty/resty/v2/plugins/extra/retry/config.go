package retry

import (
	"github.com/xgodev/boost/factory/contrib/go-resty/resty/v2"
	"time"

	"github.com/xgodev/boost/wrapper/config"
)

const (
	root        = resty.PluginsRoot + ".retry"
	enabled     = ".enabled"
	count       = ".count"
	waitTime    = ".waitTime"
	maxWaitTime = ".maxWaitTime"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+enabled, true, "enable/disable retry")
	config.Add(path+count, 0, "defines global max http retries")
	config.Add(path+waitTime, 200*time.Millisecond, "defines global retry wait time")
	config.Add(path+maxWaitTime, 2*time.Second, "defines global max retry wait time")
}
