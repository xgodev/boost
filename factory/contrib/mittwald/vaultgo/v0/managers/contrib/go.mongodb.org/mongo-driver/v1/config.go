package mongo

import (
	vault "github.com/xgodev/boost/factory/contrib/mittwald/vaultgo/v0"
)

const (
	root = vault.ManagersRoot + ".mongo"
)

func init() {
	vault.ManagerConfigAdd(root)
}
