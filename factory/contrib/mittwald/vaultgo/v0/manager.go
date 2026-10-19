package vault

import "context"

// Manager keeps a client configured with the credentials a vault secret
// holds. The pool reads the secret, maps its keys with the manager's options
// and hands the result to Configure; when the secret lease ends, it calls
// Close and configures the manager again.
type Manager interface {
	Options() *ManagerOptions
	Close(context.Context) error
	Configure(context.Context, map[string]interface{}) error
}
