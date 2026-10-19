package gocql

import (
	"time"

	"github.com/xgodev/boost/wrapper/config"
)

// Options are the session settings, read from boost.factory.gocql.*. A
// zero value keeps the driver default, except DefaultTimestamp and
// DisableInitialHostLookup, which are always applied.
type Options struct {
	Hosts                    []string
	Port                     int
	DC                       string `config:"dc"`
	Username                 string
	Password                 string
	CQLVersion               string `config:"CQLVersion"`
	ProtoVersion             int
	Timeout                  time.Duration
	ConnectTimeout           time.Duration
	Keyspace                 string
	NumConns                 int
	Consistency              string
	SocketKeepalive          time.Duration
	MaxPreparedStmts         int
	MaxRoutingKeyInfo        int
	PageSize                 int
	DefaultTimestamp         bool
	ReconnectInterval        time.Duration
	MaxWaitSchemaAgreement   time.Duration
	DisableInitialHostLookup bool
	WriteCoalesceWaitTime    time.Duration
}

// NewOptions returns options from config file or environment vars.
func NewOptions() (*Options, error) {
	return config.NewOptionsWithPath[Options](root)
}

// NewOptionsWithPath returns the options of path, over the root defaults.
func NewOptionsWithPath(path string) (opts *Options, err error) {
	return config.NewOptionsWithPath[Options](root, path)
}
