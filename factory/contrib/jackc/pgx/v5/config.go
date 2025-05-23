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
	config.Add(path+connectString, "postgres://myuser:mypassword@db.example.com:5432/mydatabase?sslmode=disable&application_name=myapp", "sets database connection string")
	config.Add(path+pu, "", "sets database username")
	config.Add(path+pp, "", "sets database password")
}
