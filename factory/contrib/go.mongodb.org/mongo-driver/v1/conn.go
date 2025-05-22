package mongo

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/xgodev/boost/model/errors"
	"github.com/xgodev/boost/wrapper/log"
	"io/ioutil"
	"strings"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

// Conn represents a mongo connection.
type Conn struct {
	ClientOptions *options.ClientOptions
	Client        *mongo.Client
	Database      *mongo.Database
	Options       *Options
	Plugins       []Plugin
}

// ClientOptionsPlugin defines a mongo client options plugin signature.
type ClientOptionsPlugin func(context.Context, *options.ClientOptions) error

// ClientPlugin defines a mongo client plugin signature.
type ClientPlugin func(context.Context, *mongo.Client) error

// Plugin defines a function to process plugin.
type Plugin func(context.Context) (ClientOptionsPlugin, ClientPlugin)

// NewConn returns a new connection with default options.
func NewConn(ctx context.Context, plugins ...Plugin) (*Conn, error) {
	logger := log.FromContext(ctx)

	o, err := NewOptions()
	if err != nil {
		logger.Errorf("Failed to get default options: %v", err)
		return nil, errors.NewInternal(err, "Failed to get default options")
	}

	return NewConnWithOptions(ctx, o, plugins...)
}

// NewConnWithConfigPath returns a new connection with options from config path.
func NewConnWithConfigPath(ctx context.Context, path string, plugins ...Plugin) (*Conn, error) {
	opts, err := NewOptionsWithPath(path)
	if err != nil {
		return nil, errors.NewInternal(err, "Failed to get options from config path")
	}
	return NewConnWithOptions(ctx, opts, plugins...)
}

// NewConnWithOptions returns a new connection with options from config path.
func NewConnWithOptions(ctx context.Context, o *Options, plugins ...Plugin) (conn *Conn, err error) {
	logger := log.FromContext(ctx)

	var clientOptionsPlugins []ClientOptionsPlugin
	var clientPlugins []ClientPlugin

	for _, plugin := range plugins {
		clientOptionsPlugin, clientPlugin := plugin(ctx)
		if clientOptionsPlugin != nil {
			clientOptionsPlugins = append(clientOptionsPlugins, clientOptionsPlugin)
		}
		if clientPlugin != nil {
			clientPlugins = append(clientPlugins, clientPlugin)
		}
	}

	co, err := clientOptions(ctx, o)
	if err != nil {
		logger.Errorf("Failed to create client options: %v", err)
		return nil, errors.NewInternal(err, "Failed to create client options")
	}

	for _, clientOptionsPlugin := range clientOptionsPlugins {
		if err := clientOptionsPlugin(ctx, co); err != nil {
			logger.Errorf("Failed to apply client options plugin: %v", err)
			return nil, errors.NewInternal(err, "Failed to apply client options plugin")
		}
	}

	var client *mongo.Client
	var database *mongo.Database

	client, database, err = newClient(ctx, co)
	if err != nil {
		return nil, errors.NewInternal(err, "Failed to create MongoDB client")
	}

	for _, clientPlugin := range clientPlugins {
		if err := clientPlugin(ctx, client); err != nil {
			logger.Errorf("Failed to apply client plugin: %v", err)
			return nil, errors.NewInternal(err, "Failed to apply client plugin")
		}
	}

	conn = &Conn{
		ClientOptions: co,
		Client:        client,
		Database:      database,
		Plugins:       plugins,
		Options:       o,
	}

	return conn, err
}

func newClient(ctx context.Context, co *options.ClientOptions) (client *mongo.Client, database *mongo.Database, err error) {
	logger := log.FromContext(ctx)

	client, err = mongo.Connect(ctx, co)
	if err != nil {
		return nil, nil, errors.NewInternal(err, "Failed to connect to MongoDB")
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, nil, errors.NewInternal(err, "Failed to ping MongoDB server")
	}

	var connFields *connstring.ConnString
	connFields, err = connstring.Parse(co.GetURI())
	if err != nil {
		return nil, nil, errors.NewInternal(err, "Failed to parse MongoDB connection string")
	}

	database = client.Database(connFields.Database)
	logger.Infof("Connected to MongoDB server: %v", strings.Join(connFields.Hosts, ","))

	return client, database, err
}

func clientOptions(ctx context.Context, o *Options) (*options.ClientOptions, error) {
	logger := log.FromContext(ctx)

	// Use the ToClientOptions method from Options struct to get the base configuration
	clientOptions := o.ToClientOptions()

	// Add command and pool monitoring
	clientOptions.SetMonitor(&event.CommandMonitor{
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			logger.Debugf("mongodb cmd - %v %s %s %v", startedEvent.ConnectionID, startedEvent.CommandName, startedEvent.DatabaseName, startedEvent.RequestID)
		},
		Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {
			logger.Debugf("mongodb cmd - %v %s %vus %v", succeededEvent.ConnectionID, succeededEvent.CommandName, succeededEvent.DurationNanos, succeededEvent.RequestID)
		},
		Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {
			logger.Errorf("mongodb cmd - %v %s %s %v", failedEvent.ConnectionID, failedEvent.CommandName, failedEvent.Failure, failedEvent.RequestID)
		},
	})
	clientOptions.SetPoolMonitor(&event.PoolMonitor{
		Event: func(poolEvent *event.PoolEvent) {
			logger.Debugf("mongodb conn pool - %v %s %s %s", poolEvent.ConnectionID, poolEvent.Type, poolEvent.Reason, poolEvent.Address)
		},
	})

	// Configure TLS if enabled
	if o.TLS != nil && *o.TLS {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: o.TLSInsecure != nil && *o.TLSInsecure,
		}

		// Load CA certificate if provided
		if o.TLSCAFile != "" {
			caCert, err := ioutil.ReadFile(o.TLSCAFile)
			if err != nil {
				return nil, errors.NewInternal(err, "Failed to read TLS CA file")
			}
			
			caCertPool := x509.NewCertPool()
			if !caCertPool.AppendCertsFromPEM(caCert) {
				return nil, errors.NewInternal(nil, "Failed to append CA certificate to pool")
			}
			
			tlsConfig.RootCAs = caCertPool
		}

		// Load client certificate if provided
		if o.TLSCertificateKeyFile != "" {
			cert, err := tls.LoadX509KeyPair(o.TLSCertificateKeyFile, o.TLSCertificateKeyFile)
			if err != nil {
				return nil, errors.NewInternal(err, "Failed to load TLS certificate/key pair")
			}
			
			tlsConfig.Certificates = []tls.Certificate{cert}
		}

		clientOptions.SetTLSConfig(tlsConfig)
	}

	return clientOptions, nil
}
