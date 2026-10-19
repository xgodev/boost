package health

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ginats "github.com/xgodev/boost/factory/contrib/nats-io/nats.go/v1"
)

var _ ginats.Plugin = (&Health{}).Plugin

func TestHealthPluginRegistersOnTheConnection(t *testing.T) {
	h := NewHealthWithOptions(&Options{Name: "nats", Enabled: true})

	connPlugin, optsPlugin := h.Plugin(context.Background())
	assert.Nil(t, optsPlugin, "health does not change the connection options")
	require.NotNil(t, connPlugin)
}
