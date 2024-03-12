package redis

import (
	vault "github.com/xgodev/boost/factory/mittwald/vaultgo.v0"
)

const (
	root = vault.ManagersRoot + ".redis"
)

func init() {
	vault.ManagerConfigAdd(root)
}
