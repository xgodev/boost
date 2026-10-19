package mongo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xgodev/boost/factory/contrib/go.mongodb.org/mongo-driver/v1"
	vault "github.com/xgodev/boost/factory/contrib/mittwald/vaultgo/v0"
)

func TestConfigureNeedsUsernameAndPassword(t *testing.T) {
	m := NewManagerWithOptions(&vault.ManagerOptions{}, &mongo.Options{})

	assert.Error(t, m.Configure(context.Background(), map[string]interface{}{"password": "p"}))
	assert.Error(t, m.Configure(context.Background(), map[string]interface{}{"username": "u"}))
	_, ok := m.Client()
	assert.False(t, ok)
}
