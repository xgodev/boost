package gocql

import (
	vault "github.com/xgodev/boost/factory/contrib/mittwald/vaultgo/v0"
)

const (
	root = vault.ManagersRoot + ".gocql"
)

func init() {
	vault.ManagerConfigAdd(root)
}
