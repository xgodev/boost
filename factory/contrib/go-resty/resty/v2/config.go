package resty

import (
	"github.com/xgodev/boost/wrapper/config"
	"time"
)

const (
	root                           = "boost.factory.resty"
	host                           = ".host"
	debug                          = ".debug"
	accept                         = ".accept"
	headers                        = ".headers"
	queryParams                    = ".queryParams"
	authorization                  = ".authorization"
	closeConnection                = ".closeConnection"
	connectionTimeout              = ".connectionTimeout"
	keepAlive                      = ".keepAlive"
	fallbackDelay                  = ".fallbackDelay"
	requestTimeout                 = ".requestTimeout"
	transportDisableCompression    = ".transport.disableCompression"
	transportDisableKeepAlives     = ".transport.disableKeepAlives"
	transportMaxIdleConnsPerHost   = ".transport.maxIdleConnsPerHost"
	transportResponseHeaderTimeout = ".transport.responseHeaderTimeout"
	transportForceAttemptHTTP2     = ".transport.forceAttemptHTTP2"
	transportMaxIdleConns          = ".transport.maxIdleConns"
	transportMaxConnsPerHost       = ".transport.maxConnsPerHost"
	transportIdleConnTimeout       = ".transport.idleConnTimeout"
	transportTLSHandshakeTimeout   = ".transport.TLSHandshakeTimeout"
	transportExpectContinueTimeout = ".transport.expectContinueTimeout"
	PluginsRoot                    = root + ".plugins"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+host, "http://localhost", "defines host request")
	config.Add(path+debug, false, "defines debug request")
	config.Add(path+accept, "application/json", "defines accept request")
	config.Add(path+headers, map[string]string{}, "defines headers request")
	config.Add(path+queryParams, map[string]string{}, "defines queryParams request")
	config.Add(path+authorization, "", "defines authorization request")
	config.Add(path+closeConnection, false, "defines http close connection")
	config.Add(path+connectionTimeout, 3*time.Minute, "defines http connection timeout")
	config.Add(path+keepAlive, 30*time.Second, "defines http keepalive")
	config.Add(path+fallbackDelay, 0*time.Millisecond, "defines fallbackDelay")
	config.Add(path+requestTimeout, 2*time.Second, "defines http request timeout")
	config.Add(path+transportDisableCompression, false, "enabled/disable transport compression")
	config.Add(path+transportDisableKeepAlives, false, "enabled/disable transport keep alives")
	config.Add(path+transportMaxIdleConnsPerHost, 100, "define transport max idle conns per host")
	config.Add(path+transportResponseHeaderTimeout, 0*time.Second, "define transport response header timeout")
	config.Add(path+transportForceAttemptHTTP2, true, "define transport force attempt http2")
	config.Add(path+transportMaxIdleConns, 100, "define transport max idle conns")
	config.Add(path+transportMaxConnsPerHost, 100, "define transport max conns per host")
	config.Add(path+transportIdleConnTimeout, 90*time.Second, "define transport idle conn timeout")
	config.Add(path+transportTLSHandshakeTimeout, 10*time.Second, "define transport TLS handshake timeout")
	config.Add(path+transportExpectContinueTimeout, 1*time.Second, "define transport expect continue timeout")
}
