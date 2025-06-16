package sql

import (
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root            = "boost.factory.sql"
	connMaxLifetime = ".connMaxLifetime"
	connMaxIdletime = ".connMaxIdletime"
	maxIdleConns    = ".maxIdleConns"
	maxOpenConns    = ".maxOpenConns"
	PluginsRoot     = root + ".plugins"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+connMaxLifetime, 0, "sets the maximum amount of time a connection may be reused. If d <= 0, connections are reused forever")
	config.Add(path+connMaxIdletime, 0, "sets the maximum amount of time a connection may be idle. If d <= 0, connections are reused forever")
	config.Add(path+maxIdleConns, 2, "sets the maximum number of idle connections in the pool")
	config.Add(path+maxOpenConns, 0, "sets the maximum number of open connections to the database")
}
