package api

import "github.com/xgodev/boost/wrapper/config"

// ConfigAdd registers shared API-level GCP configuration keys under the given root path.
func ConfigAdd(path string) {
	config.Add(path+".projectId", "", "GCP project ID")
	config.Add(path+".credentials.file", "credentials.json", "path to credentials file")
	config.Add(path+".credentials.json", "", "GCP credentials JSON")
	config.Add(path+".endpoint", "", "override API endpoint")
	config.Add(path+".useEmulator", "false", "use emulator if true")
	config.Add(path+".emulatorHost", "", "emulator host address")
	config.Add(path+".userAgent", "", "custom User-Agent header")
	config.Add(path+".scopes", "", "comma-separated OAuth scopes")
	config.Add(path+".timeout", "30s", "default RPC timeout")
	config.Add(path+".proxy", "", "HTTP(S) proxy URL")
	config.Add(path+".retry.maxAttempts", "4", "max retry attempts")
	config.Add(path+".retry.initialBackoff", "100ms", "initial retry backoff")
	config.Add(path+".retry.maxBackoff", "10s", "max retry backoff")
	config.Add(path+".retry.multiplier", "1.5", "retry backoff multiplier")
}
