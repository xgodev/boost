package mongo

import (
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root        = "boost.factory.mongo"
	uri         = ".uri"
	authRoot    = ".auth"
	username    = authRoot + ".username"
	password    = authRoot + ".password"
	PluginsRoot = root + ".plugins"
	
	// Pool configuration
	poolRoot         = ".pool"
	maxPoolSize      = poolRoot + ".max_size"
	minPoolSize      = poolRoot + ".min_size"
	maxConnIdleTime  = poolRoot + ".max_idle_time_ms"
	
	// Timeout configuration
	timeoutRoot             = ".timeout"
	connectTimeout          = timeoutRoot + ".connect_ms"
	socketTimeout           = timeoutRoot + ".socket_ms"
	serverSelectionTimeout  = timeoutRoot + ".server_selection_ms"
	heartbeatInterval       = timeoutRoot + ".heartbeat_ms"
	localThreshold          = timeoutRoot + ".local_threshold_ms"
	maxConnecting           = timeoutRoot + ".max_connecting"
	disableOCSPEndpointCheck = timeoutRoot + ".disable_ocsp_endpoint_check"
	
	// TLS/SSL configuration
	tlsRoot                 = ".tls"
	tlsEnabled              = tlsRoot + ".enabled"
	tlsInsecure             = tlsRoot + ".insecure"
	tlsCertificateKeyFile   = tlsRoot + ".certificate_key_file"
	tlsCertificatePassword  = tlsRoot + ".certificate_password"
	tlsCAFile               = tlsRoot + ".ca_file"
	
	// Read/Write concern
	readConcernRoot         = ".read_concern"
	readConcern             = readConcernRoot + ".level"
	readPreferenceRoot      = ".read_preference"
	readPreference          = readPreferenceRoot + ".mode"
	readPreferenceMaxStaleness = readPreferenceRoot + ".max_staleness_ms"
	
	writeConcernRoot        = ".write_concern"
	writeConcern            = writeConcernRoot + ".level"
	writeConcernW           = writeConcernRoot + ".w"
	writeConcernJ           = writeConcernRoot + ".j"
	writeConcernWTimeout    = writeConcernRoot + ".wtimeout_ms"
	
	// Compression
	compressionRoot         = ".compression"
	compressors             = compressionRoot + ".compressors"
	zlibLevel               = compressionRoot + ".zlib_level"
	zstdLevel               = compressionRoot + ".zstd_level"
	
	// Retry configuration
	retryRoot               = ".retry"
	retryReads              = retryRoot + ".reads"
	retryWrites             = retryRoot + ".writes"
	
	// Replica set configuration
	replicaSetRoot          = ".replica_set"
	replicaSet              = replicaSetRoot + ".name"
	direct                  = replicaSetRoot + ".direct"
	
	// Server API configuration
	serverAPIRoot           = ".server_api"
	serverAPIVersion        = serverAPIRoot + ".version"
	serverAPIStrict         = serverAPIRoot + ".strict"
	serverAPIDeprecationErrors = serverAPIRoot + ".deprecation_errors"
	
	// Miscellaneous
	appName                 = ".app_name"
	loadBalanced            = ".load_balanced"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	// Basic configuration
	config.Add(path+uri, "mongodb://localhost:27017/temp", "MongoDB connection URI")
	config.Add(path+username, "", "MongoDB username", config.WithHide())
	config.Add(path+password, "", "MongoDB password", config.WithHide())
	
	// Pool configuration
	config.Add(path+maxPoolSize, 100, "Maximum number of connections in the connection pool")
	config.Add(path+minPoolSize, 0, "Minimum number of connections in the connection pool")
	config.Add(path+maxConnIdleTime, 0, "Maximum idle time for a pooled connection in milliseconds (0 = no limit)")
	
	// Timeout configuration
	config.Add(path+connectTimeout, 30000, "Timeout for initial connection in milliseconds")
	config.Add(path+socketTimeout, 0, "Timeout for socket operations in milliseconds (0 = no timeout)")
	config.Add(path+serverSelectionTimeout, 30000, "Timeout for server selection in milliseconds")
	config.Add(path+heartbeatInterval, 10000, "Interval between server monitoring checks in milliseconds")
	config.Add(path+localThreshold, 15, "The maximum latency difference between the fastest server and acceptable servers")
	config.Add(path+maxConnecting, 2, "Maximum number of concurrent connection attempts")
	config.Add(path+disableOCSPEndpointCheck, false, "Disable OCSP endpoint checking")
	
	// TLS/SSL configuration
	config.Add(path+tlsEnabled, false, "Enable TLS/SSL for MongoDB connection")
	config.Add(path+tlsInsecure, false, "Allow insecure TLS/SSL connections (skip verification)")
	config.Add(path+tlsCertificateKeyFile, "", "Path to the client certificate and private key file")
	config.Add(path+tlsCertificatePassword, "", "Password for the client certificate private key", config.WithHide())
	config.Add(path+tlsCAFile, "", "Path to the CA certificate file")
	
	// Read/Write concern
	config.Add(path+readConcern, "", "Read concern level (available, local, majority, linearizable, snapshot)")
	config.Add(path+readPreference, "primary", "Read preference mode (primary, primaryPreferred, secondary, secondaryPreferred, nearest)")
	config.Add(path+readPreferenceMaxStaleness, 90000, "Maximum staleness for secondary reads in milliseconds")
	
	config.Add(path+writeConcern, "", "Write concern level (majority, etc.)")
	config.Add(path+writeConcernW, 1, "Write concern w value (number of nodes that must acknowledge writes)")
	config.Add(path+writeConcernJ, false, "Write concern j value (whether writes should be journaled)")
	config.Add(path+writeConcernWTimeout, 0, "Write concern timeout in milliseconds (0 = no timeout)")
	
	// Compression
	config.Add(path+compressors, []string{}, "Compression algorithms to use (snappy, zlib, zstd)")
	config.Add(path+zlibLevel, 6, "Compression level for zlib (0-9)")
	config.Add(path+zstdLevel, 6, "Compression level for zstd (1-20)")
	
	// Retry configuration
	config.Add(path+retryReads, true, "Enable retryable reads")
	config.Add(path+retryWrites, true, "Enable retryable writes")
	
	// Replica set configuration
	config.Add(path+replicaSet, "", "Replica set name")
	config.Add(path+direct, false, "Use direct connection (bypass mongos)")
	
	// Server API configuration
	config.Add(path+serverAPIVersion, "", "Server API version (1)")
	config.Add(path+serverAPIStrict, false, "Enable strict server API mode")
	config.Add(path+serverAPIDeprecationErrors, false, "Treat deprecated server API errors as errors")
	
	// Miscellaneous
	config.Add(path+appName, "", "Application name for MongoDB logs and profiling")
	config.Add(path+loadBalanced, false, "Enable load balanced mode")
}
