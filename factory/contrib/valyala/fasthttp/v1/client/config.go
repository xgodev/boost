package client

import (
	"time"

	"github.com/xgodev/boost/wrapper/config"
)

const (
	root                          = "boost.factory.fasthttp.client"
	name                          = ".name"
	noDefaultUserAgentHeader      = ".noDefaultUserAgentHeader"
	maxConnsPerHost               = ".maxConnsPerHost"
	maxConnWaitTimeout            = ".maxConnWaitTimeout"
	readBufferSize                = ".readBufferSize"
	writeBufferSize               = ".writeBufferSize"
	readTimeout                   = ".readTimeout"
	writeTimeout                  = ".writeTimeout"
	maxIdleConnDuration           = ".maxIdleConnDuration"
	maxConnDuration               = ".maxConnDuration"
	disableHeaderNamesNormalizing = ".disableHeaderNamesNormalizing"
	dialDualStack                 = ".dialDualStack"
	maxResponseBodySize           = ".maxResponseBodySize"
	maxIdemponentCallAttempts     = ".maxIdemponentCallAttempts"
)

func init() {
	ConfigAdd(root)
}

// ConfigAdd registers the fasthttp client config keys under path.
func ConfigAdd(path string) {
	config.Add(path+name, "fasthttp-client", "used in User-Agent request header")
	config.Add(path+noDefaultUserAgentHeader, false, "excludes the User-Agent header from the request")
	config.Add(path+maxConnsPerHost, 512, "maximum number of concurrent connections per host")
	config.Add(path+readBufferSize, 0, "per-connection buffer size for responses' reading, 0 is the fasthttp default")
	config.Add(path+writeBufferSize, 0, "per-connection buffer size for requests' writing, 0 is the fasthttp default")
	config.Add(path+maxConnWaitTimeout, 0*time.Second, "maximum amount of time to wait for a connection to be free, 0s fails at once")
	config.Add(path+readTimeout, 0*time.Second, "maximum duration for full response reading (including body), 0s is unlimited")
	config.Add(path+writeTimeout, 0*time.Second, "maximum duration for full request writing (including body), 0s is unlimited")
	config.Add(path+maxIdleConnDuration, 10*time.Second, "idle keep-alive connections are closed after this duration")
	config.Add(path+maxConnDuration, 0*time.Second, "keep-alive connections are closed after this duration, 0s is unlimited")
	config.Add(path+disableHeaderNamesNormalizing, false, "header names are passed as-is without normalization")
	config.Add(path+dialDualStack, false, "attempt to connect to both ipv4 and ipv6 addresses")
	config.Add(path+maxResponseBodySize, 52428800, "maximum response body size in bytes")
	config.Add(path+maxIdemponentCallAttempts, 5, "maximum number of attempts for idempotent calls")
}
