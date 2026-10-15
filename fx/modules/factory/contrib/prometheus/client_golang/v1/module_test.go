package prometheus

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"go.uber.org/fx"
)

func TestModuleGraphIsComplete(t *testing.T) {
	boost.Start()

	t.Run("PusherModule", func(t *testing.T) {
		require.NoError(t, fx.ValidateApp(PusherModule()))
	})
}
