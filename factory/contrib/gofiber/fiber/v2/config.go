package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root                      = "boost.factory.fiber"
	port                      = ".port"
	typ                       = ".type"
	configRoot                = ".config"
	prefork                   = configRoot + ".prefork"
	serverHeader              = configRoot + ".serverHeader"
	strictRouting             = configRoot + ".strictRouting"
	caseSensitive             = configRoot + ".caseSensitive"
	immutable                 = configRoot + ".immutable"
	unescapePath              = configRoot + ".unescapePath"
	etag                      = configRoot + ".etag"
	bodyLimit                 = configRoot + ".bodyLimit"
	concurrency               = configRoot + ".concurrency"
	readTimeout               = configRoot + ".readTimeout"
	writeTimeout              = configRoot + ".writeTimeout"
	idleTimeout               = configRoot + ".idleTimeout"
	readBufferSize            = configRoot + ".readBufferSize"
	writeBufferSize           = configRoot + ".writeBufferSize"
	compressedFileSuffix      = configRoot + ".compressedFileSuffix"
	proxyHeader               = configRoot + ".proxyHeader"
	getOnly                   = configRoot + ".getOnly"
	disableKeepalive          = configRoot + ".disableKeepalive"
	disableDefaultDate        = configRoot + ".disableDefaultDate"
	disableDefaultContentType = configRoot + ".disableDefaultContentType"
	disableHeaderNormalizing  = configRoot + ".disableHeaderNormalizing"
	disableStartupMessage     = configRoot + ".disableStartupMessage"
	reduceMemoryUsage         = configRoot + ".reduceMemoryUsage"
	network                   = configRoot + ".network"
	enableTrustedProxyCheck   = configRoot + ".enableTrustedProxyCheck"
	trustedProxies            = configRoot + ".trustedProxies"
	// PluginsRoot is the config root every fiber plugin key lives under.
	PluginsRoot = root + ".plugins"
)

func init() {
	ConfigAdd(root)
}

// ConfigAdd registers the fiber server config keys under path.
func ConfigAdd(path string) {
	config.Add(path+port, 8082, "server http port")
	config.Add(path+typ, "REST", "defines the server type REST/HTML, the error handler answers JSON for REST")
	config.Add(path+prefork, false, "spawns multiple go processes listening on the same port")
	config.Add(path+serverHeader, "", "value of the Server http header")
	config.Add(path+strictRouting, false, "treats /foo and /foo/ as different routes")
	config.Add(path+caseSensitive, false, "treats /Foo and /foo as different routes")
	config.Add(path+immutable, false, "makes the values the context returns immutable")
	config.Add(path+unescapePath, false, "converts all encoded characters in the route back before setting the path")
	config.Add(path+etag, false, "enables ETag header generation")
	config.Add(path+bodyLimit, 4*1024*1024, "max size in bytes of a request body")
	config.Add(path+concurrency, 256*1024, "max number of concurrent connections")
	config.Add(path+readTimeout, "0s", "amount of time allowed to read the full request, 0s is unlimited")
	config.Add(path+writeTimeout, "0s", "max duration before timing out writes of the response, 0s is unlimited")
	config.Add(path+idleTimeout, "0s", "max amount of time to wait for the next request with keep-alive, 0s is unlimited")
	config.Add(path+readBufferSize, 4096, "per-connection buffer size for reading requests")
	config.Add(path+writeBufferSize, 4096, "per-connection buffer size for writing responses")
	config.Add(path+compressedFileSuffix, ".fiber.gz", "suffix added to the compressed file name")
	config.Add(path+proxyHeader, "", "header the client ip is read from, as X-Forwarded-For")
	config.Add(path+getOnly, false, "rejects all non-GET requests")
	config.Add(path+disableKeepalive, false, "closes the connection after the first response")
	config.Add(path+disableDefaultDate, false, "omits the default Date header")
	config.Add(path+disableDefaultContentType, false, "omits the default Content-Type header")
	config.Add(path+disableHeaderNormalizing, false, "sends header names as they are")
	config.Add(path+disableStartupMessage, true, "hides the fiber startup banner, boost logs the start instead")
	config.Add(path+reduceMemoryUsage, false, "aggressively reduces memory usage at the cost of higher CPU usage")
	config.Add(path+network, fiber.NetworkTCP4, "known networks are tcp, tcp4 (IPv4-only), tcp6 (IPv6-only)")
	config.Add(path+enableTrustedProxyCheck, false, "only trusts the proxy headers of the trustedProxies")
	config.Add(path+trustedProxies, []string{}, "ips or ranges the proxy headers are trusted from")
}
