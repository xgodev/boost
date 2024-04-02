package godror

import (
	vault "github.com/xgodev/boost/factory/contrib/mittwald/vaultgo/v0"
)

const (
	root = vault.ManagersRoot + ".godror"
)

func init() {
	vault.ManagerConfigAdd(root)
}
