package otel

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xgodev/boost"
	"go.uber.org/fx"
)

func TestModuleGraphIsComplete(t *testing.T) {
	boost.Start()

	require.NoError(t, fx.ValidateApp(Module()))
}
