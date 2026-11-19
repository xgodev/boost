package pgx

import (
	"github.com/xgodev/boost/wrapper/config"
)

const (
	root          = "boost.factory.pgx"
	connectString = ".connectString"
	pu            = ".username"
	pp            = ".password"
	PluginsRoot   = root + ".plugins"
)

func init() {
	ConfigAdd(root)
}

func ConfigAdd(path string) {
	config.Add(path+connectString, "postgres://localhost:5432/postgres?sslmode=disable", "sets database connection string; username and password, when set, override the ones in it")
	config.Add(path+pu, "", "sets database username; when set, overrides the connection string user, PGUSER and the OS user")
	config.Add(path+pp, "", "sets database password; when set, overrides the connection string password and PGPASSWORD")
}
