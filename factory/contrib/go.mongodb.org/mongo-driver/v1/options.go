package mongo

import (
	"time"

	"github.com/xgodev/boost/wrapper/config"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// Options represents mongo client options.
type Options struct {
	Uri  string
	Auth *options.Credential

	// Pool configuration
	MaxPoolSize     *uint64
	MinPoolSize     *uint64
	MaxConnIdleTime *time.Duration

	// Timeout configuration
	ConnectTimeout          *time.Duration
	SocketTimeout           *time.Duration
	ServerSelectionTimeout  *time.Duration
	HeartbeatInterval       *time.Duration
	LocalThreshold          *time.Duration
	MaxConnecting           *uint64
	DisableOCSPEndpointCheck bool

	// TLS/SSL configuration
	TLS                     *bool
	TLSInsecure             *bool
	TLSCertificateKeyFile   string
	TLSCertificatePassword  string
	TLSCAFile               string

	// Read/Write concern
	ReadConcern             string
	ReadPreference          string
	ReadPreferenceTagSets   []map[string]string
	ReadPreferenceMaxStaleness *time.Duration
	WriteConcern            string
	WriteConcernW           int
	WriteConcernJ           *bool
	WriteConcernWTimeout    *time.Duration

	// Compression
	Compressors             []string

	// Retry configuration
	RetryReads              *bool
	RetryWrites             *bool
	RetryableWritesEnabled  *bool

	// Replica set configuration
	ReplicaSet              string
	Direct                  *bool

	// Server API configuration
	ServerAPIVersion        string
	ServerAPIStrict         *bool
	ServerAPIDeprecationErrors *bool

	// Miscellaneous
	AppName                 string
	ZlibLevel               *int
	ZstdLevel               *int
	LoadBalanced            *bool
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath unmarshals a given key path into options and returns it.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	return config.NewOptionsWithPath[Options](root, path)
}

// ToClientOptions converts the Options struct to MongoDB ClientOptions
func (o *Options) ToClientOptions() *options.ClientOptions {
	// Create empty client options first
	clientOptions := options.Client()

	// Apply authentication if provided
	if o.Auth != nil && (o.Auth.Username != "" || o.Auth.Password != "") {
		clientOptions.SetAuth(*o.Auth)
	}

	// Pool configuration
	if o.MaxPoolSize != nil {
		clientOptions.SetMaxPoolSize(*o.MaxPoolSize)
	}
	if o.MinPoolSize != nil {
		clientOptions.SetMinPoolSize(*o.MinPoolSize)
	}
	if o.MaxConnIdleTime != nil {
		clientOptions.SetMaxConnIdleTime(*o.MaxConnIdleTime)
	}

	// Timeout configuration
	if o.ConnectTimeout != nil {
		clientOptions.SetConnectTimeout(*o.ConnectTimeout)
	}
	if o.SocketTimeout != nil {
		clientOptions.SetSocketTimeout(*o.SocketTimeout)
	}
	if o.ServerSelectionTimeout != nil {
		clientOptions.SetServerSelectionTimeout(*o.ServerSelectionTimeout)
	}
	if o.HeartbeatInterval != nil {
		clientOptions.SetHeartbeatInterval(*o.HeartbeatInterval)
	}
	if o.LocalThreshold != nil {
		clientOptions.SetLocalThreshold(*o.LocalThreshold)
	}
	if o.MaxConnecting != nil {
		clientOptions.SetMaxConnecting(*o.MaxConnecting)
	}
	clientOptions.SetDisableOCSPEndpointCheck(o.DisableOCSPEndpointCheck)

	// TLS/SSL configuration
	if o.TLS != nil {
		clientOptions.SetTLSConfig(nil) // This will be configured properly in conn.go
	}

	// Read/Write concern
	if o.ReadConcern != "" {
		var rc *readconcern.ReadConcern
		switch o.ReadConcern {
		case "local":
			rc = readconcern.Local()
		case "majority":
			rc = readconcern.Majority()
		case "linearizable":
			rc = readconcern.Linearizable()
		case "available":
			rc = readconcern.Available()
		case "snapshot":
			rc = readconcern.Snapshot()
		default:
			rc = readconcern.New(readconcern.Level(o.ReadConcern))
		}
		clientOptions.SetReadConcern(rc)
	}

	// Read preference
	if o.ReadPreference != "" {
		var rp *readpref.ReadPref
		var err error
		
		mode := readpref.PrimaryMode
		switch o.ReadPreference {
		case "primary":
			mode = readpref.PrimaryMode
		case "primaryPreferred":
			mode = readpref.PrimaryPreferredMode
		case "secondary":
			mode = readpref.SecondaryMode
		case "secondaryPreferred":
			mode = readpref.SecondaryPreferredMode
		case "nearest":
			mode = readpref.NearestMode
		}
		
		// Handle tag sets using the proper approach for the driver version
		rpOpts := []readpref.Option{}
		
		// Convert tag sets to the format expected by the driver
		if len(o.ReadPreferenceTagSets) > 0 {
			// Create tag sets directly as []map[string]string
			rpOpts = append(rpOpts, readpref.WithTagSetsFromMaps(o.ReadPreferenceTagSets))
		}
		
		if o.ReadPreferenceMaxStaleness != nil {
			rpOpts = append(rpOpts, readpref.WithMaxStaleness(*o.ReadPreferenceMaxStaleness))
		}
		
		rp, err = readpref.New(mode, rpOpts...)
		if err == nil {
			clientOptions.SetReadPreference(rp)
		}
	}

	// Write concern
	if o.WriteConcern != "" || o.WriteConcernW > 0 || o.WriteConcernJ != nil || o.WriteConcernWTimeout != nil {
		wc := writeconcern.New()
		
		if o.WriteConcern != "" {
			switch o.WriteConcern {
			case "majority":
				wc = writeconcern.New(writeconcern.WMajority())
			default:
				wc = writeconcern.New(writeconcern.W(o.WriteConcern))
			}
		} else if o.WriteConcernW > 0 {
			wc = writeconcern.New(writeconcern.W(o.WriteConcernW))
		}
		
		if o.WriteConcernJ != nil {
			wc = writeconcern.New(writeconcern.J(*o.WriteConcernJ))
		}
		
		if o.WriteConcernWTimeout != nil {
			wc = writeconcern.New(writeconcern.WTimeout(*o.WriteConcernWTimeout))
		}
		
		clientOptions.SetWriteConcern(wc)
	}

	// Compression
	if len(o.Compressors) > 0 {
		clientOptions.SetCompressors(o.Compressors)
		
		// Set compressor levels if specified
		if o.ZlibLevel != nil {
			clientOptions.SetZlibLevel(*o.ZlibLevel)
		}
		if o.ZstdLevel != nil {
			clientOptions.SetZstdLevel(*o.ZstdLevel)
		}
	}

	// Retry configuration
	if o.RetryReads != nil {
		clientOptions.SetRetryReads(*o.RetryReads)
	}
	if o.RetryWrites != nil {
		clientOptions.SetRetryWrites(*o.RetryWrites)
	}
	if o.RetryableWritesEnabled != nil {
		// This is a deprecated option, but included for backward compatibility
		clientOptions.SetRetryWrites(*o.RetryableWritesEnabled)
	}

	// Replica set configuration
	if o.ReplicaSet != "" {
		clientOptions.SetReplicaSet(o.ReplicaSet)
	}
	if o.Direct != nil {
		clientOptions.SetDirect(*o.Direct)
	}

	// Server API configuration
	if o.ServerAPIVersion != "" {
		serverAPI := options.ServerAPI(options.ServerAPIVersion(o.ServerAPIVersion))
		if o.ServerAPIStrict != nil {
			serverAPI.SetStrict(*o.ServerAPIStrict)
		}
		if o.ServerAPIDeprecationErrors != nil {
			serverAPI.SetDeprecationErrors(*o.ServerAPIDeprecationErrors)
		}
		clientOptions.SetServerAPIOptions(serverAPI)
	}

	// Miscellaneous
	if o.AppName != "" {
		clientOptions.SetAppName(o.AppName)
	}
	if o.LoadBalanced != nil {
		clientOptions.SetLoadBalanced(*o.LoadBalanced)
	}

	// Apply URI last to ensure its settings have priority
	clientOptions.ApplyURI(o.Uri)

	return clientOptions
}
