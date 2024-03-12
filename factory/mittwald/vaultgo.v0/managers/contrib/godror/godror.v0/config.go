package godror

import (
	vault "github.com/xgodev/boost/factory/mittwald/vaultgo.v0"
)

const (
	root = vault.ManagersRoot + ".godror"
)

func init() {
	vault.ManagerConfigAdd(root)
}
