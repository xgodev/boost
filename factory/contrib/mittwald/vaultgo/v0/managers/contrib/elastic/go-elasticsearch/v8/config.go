package elasticsearch

import (
	vault "github.com/xgodev/boost/factory/contrib/mittwald/vaultgo/v0"
)

const (
	root = vault.ManagersRoot + ".elasticsearch"
)

func init() {
	vault.ManagerConfigAdd(root)
}
