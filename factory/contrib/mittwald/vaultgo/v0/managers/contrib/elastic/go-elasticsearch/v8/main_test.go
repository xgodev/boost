package elasticsearch

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xgodev/boost/factory/contrib/elastic/go-elasticsearch/v8"
	vault "github.com/xgodev/boost/factory/contrib/mittwald/vaultgo/v0"
)

func TestConfigureNeedsUsernameAndPassword(t *testing.T) {
	m := NewManagerWithOptions(&vault.ManagerOptions{}, &elasticsearch.Options{})

	assert.Error(t, m.Configure(context.Background(), map[string]interface{}{"password": "p"}))
	assert.Error(t, m.Configure(context.Background(), map[string]interface{}{"username": "u"}))
	_, ok := m.Client()
	assert.False(t, ok)
}
