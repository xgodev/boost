package redis

import (
	vault "github.com/xgodev/boost/factory/contrib/mittwald/vaultgo/v0"
)

const (
	root = vault.ManagersRoot + ".redis"
)

func init() {
	vault.ManagerConfigAdd(root)
}
