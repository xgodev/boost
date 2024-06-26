package aws

import (
	"github.com/xgodev/boost/wrapper/config"
	"time"
)

const (
	root                  = "boost.factory.aws"
	aki                   = ".accessKeyId"
	sak                   = ".secretAccessKey"
	region                = ".defaultRegion"
	accountNumber         = ".defaultAccountNumber"
	customEndpoint        = ".customEndpoint"
	retryerRoot           = ".retryer"
	retryerMaxAttempts    = retryerRoot + ".maxAttempts"
	retryerHasRateLimit   = retryerRoot + ".hasRateLimit"
	httpClientRoot        = ".httpClient"
	maxIdleConnPerHost    = httpClientRoot + ".maxIdleConnPerHost"
	maxIdleConn           = httpClientRoot + ".maxIdleConn"
	maxConnsPerHost       = httpClientRoot + ".maxConnsPerHost"
	idleConnTimeout       = httpClientRoot + ".idleConnTimeout"
	disableKeepAlives     = httpClientRoot + ".disableKeepAlives"
	disableCompression    = httpClientRoot + ".disableCompression"
	forceHTTP2            = httpClientRoot + ".forceHTTP2"
	tlsHandshakeTimeout   = httpClientRoot + ".TLSHandshakeTimeout"
	timeout               = httpClientRoot + ".timeout"
	dialTimeout           = httpClientRoot + ".dialTimeout"
	keepAlive             = httpClientRoot + ".keepAlive"
	expectContinueTimeout = httpClientRoot + ".expectContinueTimeout"
	dualStack             = httpClientRoot + ".dualStack"
	PluginsRoot           = root + ".plugins"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+aki, "", "defines the aws key id", config.WithHide())
	config.Add(path+sak, "", "defines the aws secret key", config.WithHide())
	config.Add(path+region, "", "defines the aws region")
	config.Add(path+accountNumber, "", "defines the aws account number")
	config.Add(path+customEndpoint, make(OptionsCustomEndpoint), "defines custom endpoints for services")
	config.Add(path+retryerMaxAttempts, 5, "defines max attempts for rate limit")
	config.Add(path+retryerHasRateLimit, true, "defines if retryer has rate limit")
	config.Add(path+maxIdleConnPerHost, 10, "http max idle connections per host")
	config.Add(path+maxIdleConn, 100, "http max idle connections")
	config.Add(path+maxConnsPerHost, 256, "http max connections per host")
	config.Add(path+idleConnTimeout, 90*time.Second, "http idle connections timeout")
	config.Add(path+disableKeepAlives, true, "http disable keep alives")
	config.Add(path+disableCompression, true, "http disable keep alives")
	config.Add(path+forceHTTP2, true, "http force http2")
	config.Add(path+tlsHandshakeTimeout, 10*time.Second, "TLS handshake timeout")
	config.Add(path+timeout, 30*time.Second, "timeout")
	config.Add(path+dialTimeout, 5*time.Second, "dial timeout")
	config.Add(path+keepAlive, 15*time.Second, "keep alive")
	config.Add(path+expectContinueTimeout, 1*time.Second, "expect continue timeout")
	config.Add(path+dualStack, true, "dual stack")
}
